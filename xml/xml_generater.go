package xml

import (
	"bytes"
	"encoding/xml"
	"github.com/yanzongzhen/DataFormatUtils/dataformat"
	"github.com/yanzongzhen/Logger/logger"
	charset2 "golang.org/x/net/html/charset"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ValueGetter func(name string) string                             //通过node name 获取要传的值
type ValueGetterByPath func(path string) string                       //通过 element path 获取要传的值
type XmlValueTraver func(name string, path string, value string) bool //通过 element path 获取要传的值

type ResultFilter func(name string, value string) bool

func ChangeXMLNodeValue(xmlTemplateStr string, getter ValueGetter) (string, error) {
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write([]byte(xmlTemplateStr))
	decoder := xml.NewDecoder(&xmlTemplate)
	var xmlRequest bytes.Buffer
	encoder := xml.NewEncoder(&xmlRequest)
	root, err := decoder.Token()
	if err != nil {
		return "", err
	}
	tokens := make([]xml.Token, 0, 10)
	nameSpace := make(map[string]string)

	currentNodeName := ""
	for t := root; err == nil; t, err = decoder.Token() {
		switch t.(type) {
		case xml.StartElement:
			token := t.(xml.StartElement)
			var attrs []xml.Attr
			if len(token.Attr) > 0 {
				attrs = make([]xml.Attr, 0, len(token.Attr))
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrSpace := attr.Name.Space
					attrValue := attr.Value
					if attrSpace == "xmlns" {
						nameSpace[attrValue] = attrName
						attrName = attrSpace + ":" + attrName
					} else {
						attrName = nameSpace[attrSpace] + ":" + attrName
					}
					attrs = append(attrs, xml.Attr{Name: xml.Name{Space: "", Local: attrName}, Value: attrValue})
				}
			}

			name := token.Name.Local
			currentNodeName = name
			space := token.Name.Space

			if v, ok := nameSpace[space]; ok {
				name = v + ":" + name
			}
			element := xml.StartElement{Name: xml.Name{Local: name}, Attr: attrs}
			tokens = append(tokens, element)
			break
		case xml.EndElement:
			currentNodeName = ""
			token := t.(xml.EndElement)
			name := token.Name.Local
			if v, ok := nameSpace[token.Name.Space]; ok {
				name = v + ":" + name
			}
			endElement := xml.EndElement{Name: xml.Name{Local: name}}
			tokens = append(tokens, endElement)
			break
		case xml.CharData:
			token := t.(xml.CharData)
			content := string([]byte(token))
			if len(content) > 0 {
				tokens = append(tokens, xml.CharData(getter(currentNodeName)))
			}
			break
		default:
			break
		}
	}

	for _, tok := range tokens {
		err := encoder.EncodeToken(tok)
		if err != nil {
			return "", err
		}
	}

	err = encoder.Flush()
	if err != nil {
		return "", err
	}
	return xmlRequest.String(), nil
}

//func TrvserXmlc

func DealXMLNode(result []byte, filter ResultFilter) error {
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write(result)
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReaderLabel("UTF-8", input)
	}
	root, err := decoder.Token()
	if err != nil {
		return err
	}
	currentNodeName := ""
	for t := root; err == nil; t, err = decoder.Token() {
		isBreak := false
		switch t.(type) {
		case xml.StartElement:
			token := t.(xml.StartElement)
			name := token.Name.Local
			currentNodeName = name
			break
		case xml.EndElement:
			currentNodeName = ""
			break
		case xml.CharData:
			token := t.(xml.CharData)
			//log.Printf("token len :%d %s\n", len(token), string(token))

			//log.Printf("after unEscapeText: %s \n", html.UnescapeString(string(token)))

			if len(token) > 0 {
				if filter(currentNodeName, string(token)) {
					isBreak = true
				}
			}
			break
		default:
			break
		}
		if isBreak {
			break
		}
	}
	return nil
}

func ChangeXMLValue(xmlTemplateStr string, getter ValueGetterByPath) (string, error) {
	//func(path string) string {
	//	for _, innerArg := range arg.BuildConfig.Arguments {
	//		if innerArg.Location != path {
	//			continue
	//		}
	//		return innerArg.GetValue(params)
	//	}
	//	return ""
	//}
	//arg.BuildConfig.Template
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write([]byte(xmlTemplateStr))
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReader(input, "UTF-8")
	}
	var xmlRequest bytes.Buffer
	encoder := xml.NewEncoder(&xmlRequest)
	root, err := decoder.Token()
	if err != nil {
		return "", err
	}
	tokens := make([]xml.Token, 0, 10)
	nameSpace := make(map[string]string)

	elementPath := make([]string, 0, 4)
	for t := root; err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:

			name := token.Name.Local
			//elementPath = append(elementPath, name)
			space := token.Name.Space
			var attrs []xml.Attr
			if len(token.Attr) > 0 {
				attrs = make([]xml.Attr, 0, len(token.Attr))
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrSpace := attr.Name.Space
					attrValue := attr.Value
					elementPath = append(elementPath, name+"["+attrName+"]")
					value := getter(strings.Join(elementPath, "."))
					if len(value) > 0 {
						attrValue = value
					}
					if attrSpace == "xmlns" {
						nameSpace[attrValue] = attrName
						attrName = attrSpace + ":" + attrName
					} else {
						if len(attrSpace) > 0 {
							attrName = nameSpace[attrSpace] + ":" + attrName
						}
					}
					attrs = append(attrs, xml.Attr{Name: xml.Name{Space: "", Local: attrName}, Value: attrValue})
					elementPath = elementPath[0 : len(elementPath)-1]
				}
				elementPath = append(elementPath, name)
			} else {
				elementPath = append(elementPath, name)
			}
			if v, ok := nameSpace[space]; ok {
				name = v + ":" + name
			}
			element := xml.StartElement{Name: xml.Name{Local: name}, Attr: attrs}
			tokens = append(tokens, element)
			break
		case xml.EndElement:
			name := token.Name.Local
			if v, ok := nameSpace[token.Name.Space]; ok {
				name = v + ":" + name
			}
			endElement := xml.EndElement{Name: xml.Name{Local: name}}
			tokens = append(tokens, endElement)
			elementPath = elementPath[0 : len(elementPath)-1]
			break
		case xml.CharData:
			content := string([]byte(token))
			if len(content) > 0 {
				value := getter(strings.Join(elementPath, "."))
				if len(value) > 0 {
					tokens = append(tokens, xml.CharData(value))
				}
			}
			break
		case xml.ProcInst:
			procInst := token.Copy()
			tokens = append(tokens, procInst)
			break
		default:
			break
		}
	}

	for _, tok := range tokens {
		err := encoder.EncodeToken(tok)
		if err != nil {
			return "", err
		}
	}

	err = encoder.Flush()
	if err != nil {
		return "", err
	}
	return xmlRequest.String(), nil

}

//const (
//	startElement = 1
//)

func TraverseXmlNew(source []byte, traver XmlValueTraver) error {
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write(source)
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReader(input, "UTF-8")
	}
	root, err := decoder.Token()
	if err != nil {
		return err
	}
	nameSpace := make(map[string]string)

	elementPath := make([]string, 0, 4)
	isBreak := false
	currentNodeName := ""

	//priviousType := -1
	//stateQueue := make([]int, 0, 10)
	arrayCount := make(map[string]int)
	endNodeName := make(map[string]int)
	for t := root; err == nil; t, err = decoder.Token() {
		if isBreak {
			break
		}
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			logger.Debugf("name:%v", name)
			currentNodeName = name

			if len(token.Attr) > 0 {
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrSpace := attr.Name.Space
					attrValue := attr.Value

					elementPath = append(elementPath, name+"("+attrName+")")
					logger.Debugf("elementPath:%v", elementPath)
					if attrSpace == "xmlns" {
						nameSpace[attrValue] = attrName
						attrName = attrSpace + ":" + attrName
						logger.Debugf("nameSpace:%v", nameSpace)
						logger.Debugf("attrName:%v", attrName)
					} else {
						if len(attrSpace) > 0 {
							attrName = nameSpace[attrSpace] + ":" + attrName
							logger.Debugf("nameSpace:%v", nameSpace)
							logger.Debugf("attrName:%v", attrName)
						}

						logger.Debug(arrayCount)
						if traver("", strings.Join(elementPath, "."), attrValue) {
							isBreak = true
							break
						}
					}
					elementPath = elementPath[0 : len(elementPath)-1]
					logger.Debugf("elementPath:%v", elementPath)
				}
				if isBreak {
					break
				}
			}
			//没有参数，直接append name
			elementPath = append(elementPath, name)
			logger.Debugf("elementPath:%v", elementPath)
			break
		case xml.EndElement:
			endNodeName[token.Name.Local] = 1
			currentNodeName = ""

			elementPath = elementPath[0 : len(elementPath)-1]
			//stateQueue = stateQueue[:len(stateQueue)-1]
			break
		case xml.CharData:
			content := string([]byte(token))
			if len(strings.TrimSpace(content)) > 0 {
				tempPath := strings.Join(elementPath[:len(elementPath)-1], ".")
				//if stateQueue[len(stateQueue)-2] == startElement {
				if _, ok := arrayCount[tempPath]; ok {
					if _, isEnd := endNodeName[elementPath[len(elementPath)-2]]; isEnd {
						arrayCount[tempPath] = arrayCount[tempPath] + 1
						delete(endNodeName, elementPath[len(elementPath)-2])
					}
				} else {
					arrayCount[tempPath] = 0
				}
				//}
				path := strings.Join(elementPath, ".")
				isTravel := false
				if strings.HasPrefix(path, tempPath) {
					isTravel = true
					realPath := ""
					tempRealPath := ""
					for pathIndex, pathElement := range elementPath {
						realPath = realPath + pathElement
						tempRealPath = tempRealPath + pathElement

						if itemCount, isContain := arrayCount[tempRealPath]; isContain {
							realPath = realPath + "[" + strconv.Itoa(itemCount) + "]"
						}
						if pathIndex < len(elementPath)-1 {
							realPath = realPath + "."
							tempRealPath = tempRealPath + "."
						}
					}
					if traver(currentNodeName, realPath, content) {
						isBreak = true
					}
					break
				}
				//}

				if !isTravel && traver(currentNodeName, strings.Join(elementPath, "."), content) {
					isBreak = true
				}
			}
			break
		default:
			logger.Debugln("default")
			break
		}
	}
	return nil
}

func TraverseXml(source []byte, traver XmlValueTraver) error {
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write(source)
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReader(input, "UTF-8")
	}
	root, err := decoder.Token()
	if err != nil {
		return err
	}
	nameSpace := make(map[string]string)

	elementPath := make([]string, 0, 4)
	isBreak := false
	currentNodeName := ""
	for t := root; err == nil; t, err = decoder.Token() {
		if isBreak {
			break
		}
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			currentNodeName = name
			if len(token.Attr) > 0 {
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrSpace := attr.Name.Space
					attrValue := attr.Value

					elementPath = append(elementPath, name+"["+attrName+"]")
					if attrSpace == "xmlns" {
						nameSpace[attrValue] = attrName
						attrName = attrSpace + ":" + attrName
					} else {
						if len(attrSpace) > 0 {
							attrName = nameSpace[attrSpace] + ":" + attrName
						}
						if traver("", strings.Join(elementPath, "."), attrValue) {
							isBreak = true
							break
						}
					}
					elementPath = elementPath[0 : len(elementPath)-1]
				}
				if isBreak {
					break
				}
				elementPath = append(elementPath, name)
			} else {
				elementPath = append(elementPath, name)
			}
			break
		case xml.EndElement:
			currentNodeName = ""
			elementPath = elementPath[0 : len(elementPath)-1]
			break
		case xml.CharData:
			content := string([]byte(token))
			if len(strings.TrimSpace(content)) > 0 {
				if traver(currentNodeName, strings.Join(elementPath, "."), content) {
					isBreak = true
				}
			}
			break
		default:
			break
		}
	}
	return nil
}

func IsXmlMatchCondition(source []byte, array []*Array, config []*dataformat.MatchConfig) bool {
	isMatch := true
	ExistMap := make(map[string]bool)
	err := TraverseXmlIterative(source, array, func(name string, path string, value string) bool {
		var filedPath string
		path = "$." + path
		if len(config) == 0 {
			//断开遍历
			return true
		}
		for i, c := range config {
			filedPath = c.FieldPath
			if strings.Contains(c.FieldPath, "?") {
				r, _ := regexp.Compile(`\[.+]`)
				filedPath = string(r.ReplaceAllFunc([]byte(c.FieldPath), func(bytes []byte) []byte {
					return []byte("")
				}))
				path = string(r.ReplaceAllFunc([]byte(path), func(i []byte) []byte {
					return []byte("")
				}))
			}
			switch c.Condition {
			case "=":
				if path == filedPath {
					if value == c.Value {
						logger.Debugf("返回值匹配成功:%v", config[i])
						config = append(config[:i], config[i+1:]...)
						return false
					} else {
						return true
					}
				}
				break
			case "!=":
				if path == filedPath {
					if value != c.Value {
						logger.Debugf("返回值匹配成功:%v", config[i])
						config = append(config[:i], config[i+1:]...)
						return false
					} else {
						return true
					}
				}
				break
			case "exist":
				if strings.Contains(path, filedPath) {
					logger.Debugf("返回值匹配成功:%v", config[i])
					config = append(config[:i], config[i+1:]...)
					return false
				}
				break
			case "!exist":
				if strings.Contains(path, filedPath) {
					ExistMap[c.FieldPath] = true
					return false
				}
				break
			default:
				return true
			}
		}
		return false
	})
	if err != nil {
		logger.Errorln("response 匹配失败")
		return false
	}
	for i, c := range config {
		switch c.Condition {
		case "!exist":
			if !ExistMap[c.FieldPath] {
				logger.Debugf("返回值匹配成功:%v", config[i])
				config = append(config[:i], config[i+1:]...)
			}
			break
		default:
			break
		}
	}
	if len(config) > 0 {
		for _, c := range config {
			logger.Debugf("返回值匹配失败:%v", c)
		}
		isMatch = false
	}
	return isMatch
}
