package json

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type ValueTraver func(path string, value interface{}) bool

type ChangeValueTraver func(path string, oldValue interface{}) (interface{}, bool)

type AppendItemTraver func(path string, value interface{}) (int, bool)

var breakError = errors.New("break")

type ScannerState struct {
	off               int
	data              []byte
	scan              scanner
	path              string
	previousState     int
	traver            ValueTraver
	changeValueTraver ChangeValueTraver
	appendItemTraver  AppendItemTraver
}

func (state *ScannerState) scanWhile(op int) int {
	var newOp int
	for {
		if state.off >= len(state.data) {
			newOp = state.scan.eof()
			state.off = len(state.data) + 1 // mark processed EOF with len+1
		} else {
			c := state.data[state.off]
			state.off++
			newOp = state.scan.step(&state.scan, c)
		}
		if newOp != op {
			break
		}
	}
	return newOp
}

func (state *ScannerState) value() error {
	switch op := state.scanWhile(scanSkipSpace); op {
	default:
		return errors.New("parse error")

	case scanBeginArray:
		state.previousState = scanBeginArray
		return state.array()
	case scanBeginObject:
		isNeedDelete := false
		if state.previousState == scanObjectKey {
			isNeedDelete = true
			state.previousState = scanBeginObject
		}
		return state.object(isNeedDelete)

	case scanBeginLiteral:
		isNeedDelete := false
		if state.previousState == scanObjectKey {
			isNeedDelete = true
			state.previousState = scanBeginLiteral
		}
		return state.literal(isNeedDelete)
	}
}

func (state *ScannerState) object(isNeedDelete bool) error {

	for {
		op := state.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			if isNeedDelete {
				state.path = deleteLastPath(state.path)
			}
			break
		}
		start := state.off - 1
		op = state.scanWhile(scanContinue)
		item := state.data[start : state.off-1]
		key, ok := unquoteBytes(item)
		state.path = addPath(state.path, string(key))
		if ok {
			//log.Printf("key:%s\n", string(key))
			if op == scanSkipSpace {
				op = state.scanWhile(scanSkipSpace)
			}
			if op != scanObjectKey {
				return errors.New("parse error 1")
			}
			state.previousState = scanObjectKey
			err := state.value()
			if err != nil {
				return err
			}
			op = state.scanWhile(scanSkipSpace)
			if op == scanEndObject {
				if isNeedDelete {
					state.path = deleteLastPath(state.path)
				}
				//log.Printf("end======== path:%s\n", state.path)
				break
			}

			if op != scanObjectValue {
				//sta.error(errPhase)
				return errors.New("parse error 2" + state.path)
			}
		} else {
			return errors.New("parse error 3")
		}
	}
	return nil
}

func (state *ScannerState) literal(isNeedDelete bool) error {

	start := state.off - 1
	op := state.scanWhile(scanContinue)

	// Scan read one byte too far; back up.
	state.off--
	state.scan.undo(op)
	var value interface{}
	valueData := state.data[start:state.off]
	//fmt.Println(string(valueData))
	if len(valueData) == 0 {
		return errors.New("bad format")
	}

	valueLength := 0
	switch c := valueData[0]; c {
	case 'n':
		value = nil
		valueLength = 4
		break
	case 't', 'f':
		value = c == 't'
		if c == 't' {
			valueLength = 4
		} else {
			valueLength = 5
		}
		break
	case '"':
		s, ok := unquoteBytes(valueData)
		if ok {
			value = string(s)
			valueLength = len(s) + 2
		} else {
			return errors.New("to string error")
		}
		break
	default:
		var err error
		value, err = convertNumber(string(valueData))
		if err != nil {
			return err
		}
		break
	}
	//fmt.Printf("path:%s value:%v\n", d.path, value)
	//log.Printf("start:%d current off:%d \n", start, state.off)
	if state.traver != nil {
		if state.traver(state.path, value) {
			return breakError
		}
	} else if state.changeValueTraver != nil {
		newValue, isNeedChange := state.changeValueTraver(state.path, value)

		if isNeedChange {
			strRealValue := ""
			switch realValue := newValue.(type) {
			case float64, bool, int, int64:
				strRealValue = fmt.Sprintf("%v", realValue)
				break
			case string:
				//realValue = strings.Replace(realValue, "\r\n", "\n", -1)
				//realValue = html.UnescapeString(realValue)
				strRealValue = getJsonString(realValue, true)
				break

			default:
				if newValue == nil {
					strRealValue = "null"
				}
				break
			}
			//var newData []byte
			newData := make([]byte, 0, len(state.data)+len([]byte(strRealValue))-valueLength)
			newData = append(newData, state.data[0:start]...)
			newData = append(newData, []byte(strRealValue)...)
			newData = append(newData, state.data[state.off:]...)
			state.data = newData
			//if len(strRealValue) > 1000 {
			//	ioutil.WriteFile("./111.txt", []byte(strRealValue), os.ModePerm)
			//	ioutil.WriteFile("./222.json", newData, os.ModePerm)
			//}
			state.off = start + len(strRealValue)
		}
	}

	if isNeedDelete {
		state.path = deleteLastPath(state.path)
	}

	return nil
}

func getJsonString(s string, escapeHTML bool) string {
	buf := bytes.Buffer{}
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}

func (state *ScannerState) array() error {

	i := 0
	isNeedCheckArray := true
	arrayType := 0
	for {
		op := state.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			state.path = deleteLastPath(state.path)
			break
		}
		state.off--
		state.scan.undo(op)
		if isNeedCheckArray {
			nextOp := state.scanWhile(scanSkipSpace)
			if nextOp == scanBeginLiteral {
				arrayType = 1
				if strings.HasSuffix(state.path, "]") {
					arrayType = 2
				}
			} else if nextOp == scanBeginObject {
				if strings.HasSuffix(state.path, "]") {
					arrayType = 3
					state.path = state.path + "."
				} else {
					if state.path == "$" {
						state.path = state.path + "."
					}
				}

			} else if nextOp == scanBeginArray {
				//logger.Debugln(state.path)
				//return errors.New("bad format")
				if state.path == "$" {
					state.path = state.path + "."
				}

			} else {
				return errors.New("bad format")
			}
			state.off--
			state.scan.undo(nextOp)
			isNeedCheckArray = false
		}

		switch arrayType {
		case 0:
			state.path = addIndexToPath(state.path, i, true)
			break
		case 1:
			state.path = state.path + "."
			state.path = addIndexToPath(state.path, i, true)
			break
		case 2:
			state.path = state.path + "."
			state.path = addIndexToPath(state.path, i, false)
			break
		case 3:
			state.path = addIndexToPath(state.path, i, true)
			break
		}
		err := state.value()
		if err != nil {
			return err
		}
		i++
		op = state.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			if arrayType == 1 || arrayType == 3 || arrayType == 0 {
				state.path = deleteLastPath(state.path)
			}
			break
		}

		if op != scanArrayValue {
			return errors.New("parse array error")
		}
	}
	return nil
}
func addIndexToPath(path string, i int, isReplace bool) string {
	if isReplace {
		if !strings.HasSuffix(path, "]") {
			return path + "[" + strconv.Itoa(i) + "]"
		} else {
			res := IndexMatcher.FindAllStringIndex(path, -1)
			if len(res) > 0 {
				return path[:res[len(res)-1][0]] + "[" + strconv.Itoa(i) + "]"
			} else {
				return path + "[" + strconv.Itoa(i) + "]"
			}
		}
	} else {
		return path + "[" + strconv.Itoa(i) + "]"
	}
}

func addPath(path string, newPath string) string {
	return path + "." + newPath
}

func deleteLastPath(path string) string {
	return path[:strings.LastIndex(path, ".")]
}

func TravelJsonData(jsonData []byte, traver ValueTraver) error {
	//type ValueTraver func(path string, value interface{}) bool
	state := ScannerState{}
	state.path = "$"
	state.data = jsonData
	state.scan.reset()
	state.traver = traver

	err := state.value()
	if err == breakError {
		return nil
	}
	return err
}
func ChangeValueTravel(jsonData []byte, traver ChangeValueTraver) ([]byte, error) {
	state := ScannerState{}
	state.path = "$"
	state.data = jsonData
	state.scan.reset()
	state.changeValueTraver = traver

	err := state.value()
	if err == breakError {
		return state.data, nil
	}
	return state.data, nil
}

func unquote(s []byte) (t string, ok bool) {
	s, ok = unquoteBytes(s)
	t = string(s)
	return
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

			// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

			// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

			// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

func convertNumber(s string) (interface{}, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("not number ")
	}
	return f, nil
}

var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var hex = "0123456789abcdef"
