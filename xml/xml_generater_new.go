package xml

import (
	"bytes"
	"encoding/xml"
	"github.com/yanzongzhen/Logger/logger"
	charset2 "golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Count struct {
	Num      int
	Position int
	RootName string
}

type Attribute struct {
	Name  string
	Value string
}

type Token struct {
	Root       string
	Content    string
	Name       string
	Attributes []Attribute
}

func preCountToken(decoder *xml.Decoder, tt xml.Token, countMap map[int]map[string][]Count, rootName string, depth int) {
	switch token := tt.(type) {
	case xml.StartElement:
		nameSpace := make(map[string]string)
		name := token.Name.Local
		space := token.Name.Space
		if len(space) > 0 {
			for _, attr := range token.Attr {
				space := attr.Name.Space
				name := attr.Name.Local
				value := attr.Value
				if space == "xmlns" {
					nameSpace[value] = name
					break
				}
			}
			if _, ok := nameSpace[space]; ok {
				name = nameSpace[space] + ":" + name
			}
		}
		depth++
		if _, ok := countMap[depth]; !ok {
			countMap[depth] = make(map[string][]Count)
		}
		if _, ok := countMap[depth][name]; !ok {
			countMap[depth][name] = make([]Count, 0)
		}
		count := Count{0, 0, rootName}
		countMap[depth][name] = append(countMap[depth][name], count)
		index := len(countMap[depth][name])
		rootName = name + "[" + strconv.Itoa(index-1) + "]"
		break
	default:
		t, err := decoder.Token()
		if err != nil {
			return
		}
		preCountToken(decoder, t, countMap, rootName, depth)
		break
	}
	var t xml.Token
	var err error
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch t.(type) {
		case xml.StartElement:
			preCountToken(decoder, t, countMap, rootName, depth)
		case xml.EndElement:
			return
		case xml.CharData:
			break
		default:
			break
		}
		if err != nil {
			logger.Errorln(err)
		}
	}
}

func parseXMLToken(decoder *xml.Decoder, tt xml.Token, rootName string, resultMap map[int][]*Token, tokenPath []string, tokenMap map[string]*Token, depth int, countMap map[int]map[string][]Count) {
	switch token := tt.(type) {
	case xml.StartElement:
		root := &Token{}
		nameSpace := make(map[string]string)
		name := token.Name.Local
		space := token.Name.Space
		if len(space) > 0 {
			for _, attr := range token.Attr {
				space := attr.Name.Space
				name := attr.Name.Local
				value := attr.Value
				if space == "xmlns" {
					nameSpace[value] = name
					break
				}
			}
			if _, ok := nameSpace[space]; ok {
				name = nameSpace[space] + ":" + name
			}
		}

		root.Name = name
		if len(token.Attr) > 0 {
			for i, attr := range token.Attr {
				attrName := attr.Name.Local
				attrSpace := attr.Name.Space
				attrValue := attr.Value
				if attr.Name.Local == "html" {
					_ = ioutil.WriteFile("./attrvalue", []byte(attr.Value), os.ModePerm)
				}
				for ; i < len(token.Attr); i++ {
					name := token.Attr[i].Name.Local
					space := token.Attr[i].Name.Space
					value := token.Attr[i].Value
					if space == "xmlns" {
						nameSpace[value] = name
						break
					}
				}
				if attrSpace == "xmlns" {
					attrName = attrSpace + ":" + attrName
				} else {
					if len(attrSpace) > 0 {
						attrName = nameSpace[attrSpace] + ":" + attrName
					}
				}
				attr := Attribute{Name: attrName, Value: attrValue}
				root.Attributes = append(root.Attributes, attr)
			}
		}
		depth++
		tokenArray := countMap[depth][name]
		if depth == 1 {
			rootName = name
			root.Root = ""
		} else {
			root.Root = rootName
			rootNames := strings.Split(rootName, ".")
			lastName := rootNames[len(rootNames)-1]
			for i, token := range tokenArray {
				if strings.Contains(token.RootName, lastName) {
					if token.Num > 1 && token.Position != token.Num {
						rootName = rootName + "." + name + "[" + strconv.Itoa(token.Position) + "]"
						root.Name = name + "[" + strconv.Itoa(token.Position) + "]"
						token.Position++
						tokenArray[i] = token
					} else {
						rootName = rootName + "." + name
					}
				}
				//攀枝花社保参保信息不支持
				//else {
				//	rootName = rootName + "." + name
				//}
			}
		}
		tokenMap[name] = root
		tokenPath = append(tokenPath, name)
		resultMap[depth] = append(resultMap[depth], root)
		break
	default:
		//首个元素不是StartElement
		t, err := decoder.Token()
		if err != nil {
			return
		}
		parseXMLToken(decoder, t, rootName, resultMap, tokenPath, tokenMap, depth, countMap)
		break
	}
	var t xml.Token
	var err error
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			// 处理元素开始（标签）
			parseXMLToken(decoder, t, rootName, resultMap, tokenPath, tokenMap, depth, countMap)
		case xml.EndElement:
			// 处理元素结束（标签）
			tokenPath = tokenPath[0 : len(tokenPath)-1]
			return
		case xml.CharData:
			//开始和结束标签中间的文本
			content := string([]byte(token.Copy()))
			content = strings.Replace(content, "\n", "", -1)
			content = strings.TrimSpace(content)
			if len(content) > 0 {
				tokenMap[tokenPath[len(tokenPath)-1]].Content = content
			}
			break
		default:
			break
		}
	}
}

func preCount(source []byte) (map[int]map[string][]Count, error) {
	var (
		t           xml.Token
		xmlTemplate bytes.Buffer
		err         error
	)
	xmlTemplate.Write(source)
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReaderLabel("UTF-8", input)
	}
	t, err = decoder.Token()
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	countMap := make(map[int]map[string][]Count) //深度->token name->数量
	preCountToken(decoder, t, countMap, "", 0)
	//合并同类项
	for i := 1; i <= len(countMap); i++ {
		depthMap := countMap[i]
		for name, tokenArray := range depthMap {
			if len(tokenArray) > 1 {
				tokenMap := make(map[string]int)
				for _, token := range tokenArray {
					tokenMap[token.RootName]++
				}
				tokenArray = tokenArray[0:0]
				for rootName, size := range tokenMap {
					count := Count{size, 0, rootName}
					tokenArray = append(tokenArray, count)
				}
				depthMap[name] = tokenArray
			} else {
				if tokenArray[0].RootName != "" {
					tokenArray[0].Num = 1
				}
			}
		}
	}
	return countMap, nil
}

func parse(source []byte) (map[int][]*Token, error) {
	var (
		t           xml.Token
		xmlTemplate bytes.Buffer
		err         error
	)
	xmlTemplate.Write(source)
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		//logger.Debug(charset)
		return charset2.NewReaderLabel("UTF-8", input)
		//return charset2.NewReader(input, "UTF-8")
	}
	t, err = decoder.Token()
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	tokenPath := make([]string, 0)
	tokenMap := make(map[string]*Token)
	resultMap := make(map[int][]*Token)
	countMap, err := preCount(source)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	parseXMLToken(decoder, t, "", resultMap, tokenPath, tokenMap, 0, countMap)
	return resultMap, nil
}

type Array struct {
	Depth float64 `json:"depth"`
	Name  string  `json:"name"`
}

func TraverseXmlIterative(source []byte, arrayList []*Array, traver XmlValueTraver) error {
	resultMap, err := parse(source)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	//解决template中认为是数组但是只有一条数据的情况
	arrayMap := make(map[int][]string)
	for _, array := range arrayList {
		if _, ok := arrayMap[int(array.Depth)]; !ok {
			arrayMap[int(array.Depth)] = make([]string, 0)
		}
		arrayMap[int(array.Depth)] = append(arrayMap[int(array.Depth)], array.Name)
	}
travelLoop:
	for i := 1; i <= len(resultMap); i++ {
		resultArray := resultMap[i]
		for _, token := range resultArray {
			if len(arrayMap[i]) > 0 {
				for _, name := range arrayMap[i] {
					if token.Name == name {
						token.Name = token.Name + "[0]"
					}
				}
			}
			for depth := range arrayMap {
				roots := strings.Split(token.Root, ".")
				if len(roots) >= depth {
					for _, name := range arrayMap[depth] {
						if roots[depth-1] == name {
							roots[depth-1] += "[0]"
						}
						token.Root = strings.Join(roots, ".")
					}
				}
			}
			if len(token.Attributes) > 0 {
				for _, attr := range token.Attributes {
					if len(token.Root) > 0 {
						if traver("", token.Root+"."+token.Name+"."+attr.Name, attr.Value) {
							break travelLoop
						}
					} else {
						if traver("", token.Name+"."+attr.Name, attr.Value) {
							break travelLoop
						}
					}
				}
			}
			if len(token.Content) > 0 {
				if len(token.Root) > 0 {
					if traver("", token.Root+"."+token.Name, token.Content) {
						break travelLoop
					}
				} else {
					if traver("", token.Name, token.Content) {
						break travelLoop
					}
				}
			}
		}
	}
	return nil
}

func ChangeXmlValueNew(xmlTemplateStr string, getter ValueGetterByPath) (string, error) {
	var xmlTemplate bytes.Buffer
	xmlTemplate.Write([]byte(xmlTemplateStr))
	decoder := xml.NewDecoder(&xmlTemplate)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charset2.NewReaderLabel("UTF-8", input)
	}
	root, err := decoder.Token()
	if err != nil {
		return "", err
	}
	tokens := make([]xml.Token, 0, 10)
	for t := root; err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			space := token.Name.Space
			var attrs []xml.Attr
			if len(token.Attr) > 0 {
				attrs = make([]xml.Attr, 0, len(token.Attr))
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrSpace := attr.Name.Space
					attrValue := attr.Value
					if len(attrValue) > 0 && strings.HasPrefix(attrValue, "$") {
						newValue := getter(attrValue[2:])
						attrValue = newValue
					}
					attrs = append(attrs, xml.Attr{Name: xml.Name{Space: attrSpace, Local: attrName}, Value: attrValue})
				}
			}
			element := xml.StartElement{Name: xml.Name{Local: name, Space: space}, Attr: attrs}
			tokens = append(tokens, element)
			break
		case xml.EndElement:
			name := token.Name.Local
			space := token.Name.Space
			endElement := xml.EndElement{Name: xml.Name{Local: name, Space: space}}
			tokens = append(tokens, endElement)
			break
		case xml.CharData:
			content := string([]byte(token))
			if len(content) > 0 {
				if strings.HasPrefix(content, "$") {
					newContent := getter(content[2:])
					content = newContent
				}
				tokens = append(tokens, xml.CharData(content))
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
	var xmlRequest bytes.Buffer
	encoder := xml.NewEncoder(&xmlRequest)
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
