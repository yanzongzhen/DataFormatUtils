package mixformat

import (
	j "encoding/json"
	"errors"
	"fmt"
	"github.com/yanzongzhen/DataFormatUtils/json"
	"github.com/yanzongzhen/DataFormatUtils/trace"
	"github.com/yanzongzhen/DataFormatUtils/xml"
	"github.com/yanzongzhen/Logger/logger"
	"regexp"
	"strconv"
	"strings"
)

func XmlToJsonNew(xmlData []byte, arrayList []*xml.Array, template string, args map[string]string, catchError trace.CatchErrors) ([]byte, error) {
	logger.Infoln("开始进行模板匹配")
	resMap := make(map[string]interface{})
	j.Unmarshal([]byte(template), &resMap)

	//将template封装到map中
	_ = json.TravelJsonData([]byte(template), func(path string, value interface{}) bool {
		//遍历template
		if strValue, ok := value.(string); ok {
			if strings.HasPrefix(strValue, "$") { //如果value是以$开头的
				if json.HasIndexMatcher.MatchString(path) { //如果对应的path中有[n] 即存在数组
					//Path=$.data.baseInfo[0].base.collectionName
					//Value=$.response.body[0].detail[0].zjlx.string
					itemCount := getItemCountInXml(xmlData, strValue, arrayList)
					//遍历xmlData 根据路径$.response.body.detail.zjlx.string找到xml中对应的值出现的个数
					if itemCount > 1 { //如果该路径对应的值出现了多次
						pathIndex := json.IndexMatcher.FindAllStringIndex(path, -1) //[[7 15]] 查看出现数组的是哪一项
						splitIndex := pathIndex[len(pathIndex)-1][0]                //15 取数组中的最后一个数组的第二个值 取最后一个出现数组的那一项
						json.AppendItemInRes(resMap, strings.Split(path[2:splitIndex], "."), itemCount, strValue)
						//传入重复出现的次数 [data baseInfo]和resMap(template) 将最后一个出现数组的那一项包括前面的项传入
						//根据重复的次数复制出现数组的那一项并添加到resMap中
					}
				}
			}
		}
		return false
	})
	resData, err := j.Marshal(resMap)
	logger.Debug(string(resData))
	if err != nil {
		return nil, err
	}
	logger.Debugf("匹配后的template为:%v", string(resData))
	resultMap, err := getXmlResultMap(xmlData, arrayList) //resultMap[path] = value
	logger.Debugf("resData:%v", string(resData))
	//遍历xmlData,将xmlData的path和value存入resultMap中
	logger.Infof("开始向模板填充返回值")
	result, err := json.ChangeValueTravel(resData, func(path string, oldValue interface{}) (interface{}, bool) {
		//遍历resData 处理后的template

		realValuePath, ok := oldValue.(string)

		//logger.Debugf("111111111111111realValuePath:%v", realValuePath) //转换前路径
		//logger.Debugf("111111111111111path:%v", path)                   //转换后路径
		if !ok {
			logger.Errorf("%v格式匹配失败", path)
			return nil, false
		}
		if strings.HasPrefix(realValuePath, "$.") {
			pathArray := strings.Split(realValuePath, ".")
			valuePath := strings.Join(pathArray[1:len(pathArray)-1], ".") //将value最后的数据类型干掉

			var newValue interface{}
			if resultMap == nil {
				newValue, err = getValueByPathNew(xmlData, arrayList, valuePath, pathArray[len(pathArray)-1])
				if err != nil {
					return nil, false
				}
			} else {
				//todo 修改
				if strings.Contains(valuePath, "?") {
					r, _ := regexp.Compile(`\[.+]`)
					newValuePath := string(r.ReplaceAllFunc([]byte(valuePath), func(bytes []byte) []byte {
						return []byte("")
					}))
					for path, value := range resultMap {
						newPath := string(r.ReplaceAllFunc([]byte(path), func(bytes []byte) []byte {
							return []byte("")
						}))
						if newPath == newValuePath {
							newValue = value
						}
					}
				} else {
					var ok bool
					newValue, ok = resultMap[valuePath]
					if !ok {
						if json.IndexMatcher.MatchString(valuePath) {
							valuePath = json.IndexMatcher.ReplaceAllString(valuePath, "[0]")
							newValue, ok = resultMap[valuePath]
							if !ok {
								catchError(trace.FieldDisappeared, valuePath)
								logger.Errorf("返回值中不存在匹配值,匹配值:%v", valuePath)
								return "", true
							}

							newValueStr, isStr := newValue.(string)
							if isStr && len(newValueStr) <= 0 {
								catchError(trace.FieldIsEmpty, valuePath)
							}

						} else {
							return "", true
						}
					}
					newValueStr, isStr := newValue.(string)
					logger.Debugf("111111111111111newValueStr:%v", newValueStr)
					logger.Debugf("111111111111111isStr:%v", isStr)
					if isStr && len(newValueStr) <= 0 {
						logger.Debugf(valuePath)
						catchError(trace.FieldIsEmpty, valuePath)
					}
				}
			}
			newTypeValue, err := convertValueType(newValue, pathArray[len(pathArray)-1])
			if err == nil {
				return newTypeValue, true
			}
			logger.Errorf("返回值格式转换失败: %s\n", err.Error())
			return newValue, true
		} else if strings.HasPrefix(realValuePath, "#.") {
			path := strings.Split(realValuePath, ".")[1]
			logger.Debugf("path:%s", path)
			logger.Debugf("args:%v", args)
			if value, ok := args[path]; ok {
				return value, true
			}
			logger.Errorf("参数中不存在匹配值,匹配值:%v", path)
			return nil, false
		}
		return nil, false
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func getXmlResultMap(xmlData []byte, arrayList []*xml.Array) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	//todo 替换
	err := xml.TraverseXmlIterative(xmlData, arrayList, func(name string, path string, value string) bool {
		res[path] = value
		return false
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//func GetValueByPathNew(root []byte, valuePath string, valueType string) (interface{}, error) {
//	return getValueByPathNew(root, valuePath, valueType)
//}

func getValueByPathNew(root []byte, arrayList []*xml.Array, valuePath string, valueType string) (interface{}, error) {
	res := ""
	isFind := false
	//todo 替换
	err := xml.TraverseXmlIterative(root, arrayList, func(name string, path string, value string) bool {
		if valuePath == path {
			res = value
			isFind = true
			return true
		}
		return false
	})
	if !isFind {
		if json.IndexMatcher.MatchString(valuePath) {
			valuePath = json.IndexMatcher.ReplaceAllString(valuePath, "[0]")
			err = xml.TraverseXmlIterative(root, arrayList, func(name string, path string, value string) bool {
				if valuePath == path {
					res = value
					isFind = true
					return true
				}
				return false
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return convertValueType(res, valueType)
}

func getItemCountInXml(xmlData []byte, valuePath string, arrayList []*xml.Array) int {
	//Path=$.data.baseInfo[0].base.collectionName
	//Value=$.response.body[0].detail[0].zjlx.string
	count := 0

	strIndexArray := json.IndexMatcher.FindAllStringIndex(valuePath, -1)
	var tempValuePath, ralph string
	if len(strIndexArray) > 0 {
		lastIndex := strIndexArray[len(strIndexArray)-1]
		tempValuePath = valuePath[:lastIndex[0]] + "[0]" + valuePath[lastIndex[1]:]
		path := strings.Split(tempValuePath, ".")
		ralph = strings.Join(path[1:len(path)-1], ".")
	} else {
		path := strings.Split(valuePath, ".")
		ralph = strings.Join(path[1:len(path)-1], ".")
	}
	//todo 替换
	//tempValuePath := json.IndexMatcher.ReplaceAllString(valuePath, "") //Value=$.response.body.detail.zjlx.string
	err := xml.TraverseXmlIterative(xmlData, arrayList, func(name string, path string, value string) bool {
		//遍历返回值(xml)
		pathIndexArray := json.IndexMatcher.FindAllStringIndex(path, -1)
		var tempPath string

		if len(pathIndexArray) > 0 {
			lastPathIndex := pathIndexArray[len(pathIndexArray)-1]
			tempPath = path[:lastPathIndex[0]] + "[0]" + path[lastPathIndex[1]:]
		} else {
			tempPath = path
		}

		//tempPath := json.IndexMatcher.ReplaceAllString(path, "") //currentName=cxjg realPath=response.head[0].cxjg content=0
		//response.head.cxjg
		if ralph == tempPath {
			count++
		}
		return false
	})
	if err != nil {
		return -1
	}

	return count
}

//func XmlToJson(xmlData []byte, template string) ([]byte, error) {
//	templateMap := make(map[string]interface{})
//	err := j.Unmarshal([]byte(template), &templateMap)
//	if err != nil {
//		return nil, err
//	}
//
//	res := json.CopyJsonItem(templateMap)
//
//	var getValueFnc GetValue
//	getValueFnc = func(path []string, dataType string) (interface{}, error) {
//		return getValueByPath(xmlData, strings.Join(path, "."))
//	}
//	err = xml.TraverseXml(xmlData, func(name string, path string, value string) bool {
//		setTemplateValue(getValueFnc, "$."+path, "$."+path, templateMap, value, nil, res)
//		return false
//	})
//
//	if err != nil {
//		return nil, err
//	}
//	dataformat.SetResToDefault(res)
//	return j.Marshal(res)
//}

func setTemplateValue(getValueFunc GetValue, realPath string, tempPath string, template map[string]interface{}, realValue interface{}, parent interface{}, res map[string]interface{}) interface{} {
	for k, v := range template {
		switch templateValue := v.(type) {
		case string:
			if strings.HasPrefix(templateValue, "$") {

				if strings.HasPrefix(templateValue, tempPath) && templateValue[len(tempPath)] == '.' {
					parentArray, isArrayItem := parent.([]interface{}) //父亲节点是否是jsonArray
					//if hasIndexMatcher.MatchString(templateValue) { //结果路径里面包含数组
					if isArrayItem {
						isNeedNewItem := true

						for _, resArrayItem := range parentArray {
							item := resArrayItem.(map[string]interface{})
							if itemValue, ok := item[k]; ok {
								if stringValue, success := itemValue.(string); success {
									if strings.HasPrefix(stringValue, "$") {
										if strings.HasSuffix(stringValue, "e%x") {
											isNeedNewItem = false
											if len(parentArray) == 1 {
												path := strings.Split(templateValue, ".")
												res[k], _ = convertValueType(realValue, path[len(path)-2])
											}
											continue
										}
										path := strings.Split(templateValue, ".")
										res[k], _ = convertValueType(realValue, path[len(path)-1])
										isNeedNewItem = false

										for resItemKey, resItemValue := range res {
											if itemStringValue, success := resItemValue.(string); success {
												if strings.HasPrefix(itemStringValue, "$") && strings.HasSuffix(itemStringValue, "e%x") {
													//emptyValue, err := getValueByPath(root, itemStringValue[1:])
													valuePath := strings.Split(itemStringValue[1:], ".")
													emptyValue, err := getValueFunc(valuePath[:len(valuePath)-2], "string")
													if err != nil {
														res[resItemKey] = nil
													} else {
														res[resItemKey] = emptyValue
													}
												}
											}
										}
										break
									} else {
										continue
									}
								} else {
									continue
								}
							}
						}

						if isNeedNewItem {
							newItem := json.CopyJsonItem(template)
							path := strings.Split(templateValue, ".")
							if strings.HasSuffix(templateValue, "e%x") {
								newItem[k], _ = convertValueType(realValue, path[len(path)-2])
							} else {
								newItem[k], _ = convertValueType(realValue, path[len(path)-1])
							}
							parentArray = append(parentArray, newItem)
							return parentArray
						}
						return nil
					} else {
						resValue, ok := res[k].(string)
						if ok {
							if strings.HasPrefix(resValue, "$.") {
								path := strings.Split(templateValue, ".")
								res[k], _ = convertValueType(realValue, path[len(path)-1])
							}
						}

						continue
					}
				}
			} else {
				continue
			}
			break
		case []interface{}:
			resArray := res[k].([]interface{})
			setTemplateArrayValue(getValueFunc, realPath, tempPath, templateValue, realValue, k, res, resArray)
			break
		case map[string]interface{}:
			resMap := res[k].(map[string]interface{})
			setTemplateValue(getValueFunc, realPath, tempPath, templateValue, realValue, res, resMap)
			break
		}
	}
	return nil
}

func setTemplateArrayValue(getValueFunc GetValue, realPath string, tempPath string, template []interface{}, realValue interface{}, parentKey string, parent interface{}, res []interface{}) {
	for _, v := range template {
		switch value := v.(type) {
		case map[string]interface{}:
			resItem := res[len(res)-1].(map[string]interface{})
			newArray := setTemplateValue(getValueFunc, realPath, tempPath, value, realValue, res, resItem)
			if newArray != nil {
				if parentArray, ok := parent.([]interface{}); ok {
					parentArray[len(parentArray)-1] = newArray
				} else if parentMap, ok := parent.(map[string]interface{}); ok {
					parentMap[parentKey] = newArray
				}
			}
			break
		case []interface{}:
			parentValue := res[len(res)-1].([]interface{})
			setTemplateArrayValue(getValueFunc, realPath, tempPath, value, realValue, "", res, parentValue)
			break
		case string:
			if strings.HasPrefix(value, "$") {
				if strings.HasPrefix(value, tempPath) {
					path := strings.Split(value, ".")
					item, _ := convertValueType(realValue, path[len(path)-1])
					res = append(res, item)
					if parentArray, ok := parent.([]interface{}); ok {
						parentArray[len(parentArray)-1] = res
					} else if parentMap, ok := parent.(map[string]interface{}); ok {
						parentMap[parentKey] = res
					}
				}
			}
			break
		default:
			break
		}
	}
}

func ConvertValueType(value interface{}, valueType string) (interface{}, error) {
	return convertValueType(value, valueType)
}

func convertValueType(value interface{}, valueType string) (interface{}, error) {
	if valueType == "string" {
		return fmt.Sprint(value), nil
	}
	switch v := value.(type) {
	case float64:
		if valueType == "float" {
			return v, nil
		} else if valueType == "int" {
			return int(v), nil
		} else {
			logger.Errorf("类型转换失败:%v to type:%s\n", value, valueType)
			return nil, errors.New(fmt.Sprintf("%s type not favor", valueType))
		}
	case string:
		if valueType == "float" {
			return strconv.ParseFloat(v, 64)
		} else if valueType == "int" {
			fv, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, errors.New("parse float err")
			}
			return int(fv), nil
		} else {
			logger.Errorf("类型转换失败:%v to type:%s\n", value, valueType)
			return nil, errors.New(fmt.Sprintf("%s type not favor", valueType))
		}
	default:
		return v, nil
	}
}

func getValueByPath(root []byte, valuePath string) (interface{}, error) {
	res := ""
	err := xml.TraverseXml(root, func(name string, path string, value string) bool {
		if valuePath == path {
			res = value
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
