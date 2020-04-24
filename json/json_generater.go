package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yanzongzhen/DataFormatUtils/dataformat"
	"github.com/yanzongzhen/DataFormatUtils/trace"
	"github.com/yanzongzhen/Logger/logger"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var pureIndexMatcher *regexp.Regexp
var HasIndexMatcher *regexp.Regexp
var IndexMatcher *regexp.Regexp

func init() {
	pureIndexMatcher, _ = regexp.Compile(`^\[[0-9]+]$`)
	HasIndexMatcher, _ = regexp.Compile(`[^.]+?\[[0-9]+]`)
	IndexMatcher, _ = regexp.Compile(`\[[0-9]+]`)
}

//<<<<<<< Updated upstream
func TransJsonFormat(source []byte, template map[string]interface{}, catchErrors trace.CatchErrors) map[string]interface{} {
	res := copyJsonItem(template)
	root := make(map[string]interface{})
	err := json.Unmarshal(source, &root)

	logger.Debugf("1111111111111111111res:%v", res)
	logger.Debugf("1111111111111111111root:%v", root)

	//=======
	//func TransJsonFormat(body []byte, template string) []byte {
	//	res := copyJsonItem(template)
	//	root := make(map[string]interface{})
	//	err := json.Unmarshal(source, &root)
	if err != nil {
		return nil
	}

	if root == nil {
		return nil
	}
	err = TravelJsonData(source, func(path string, value interface{}) bool {

		valueStr, isStr := value.(string)
		if (isStr && len(valueStr) <= 0) || value == nil {
			if strings.HasPrefix(path, "$.") {
				catchErrors(trace.FieldIsEmpty, path[2:])
			} else {
				catchErrors(trace.FieldIsEmpty, path)
			}
		}
		//遍历返回的json 遍历时每遇到一次path和value都遍历一次template找到value对应的位置并赋值
		//log.Println(path)
		tempPath := path
		if HasIndexMatcher.MatchString(path) {
			tempPath = string(IndexMatcher.ReplaceAllString(path, "[0]"))
		}
		if value == nil {
		}
		setTemplateValue(root, path, tempPath, template, value, nil, res)
		//root 返回值map path 解析的每个路径 tempPath 将path中的[n]全部转化为[0]
		//template transTemplate的map  value path对应的值  res template的复制
		return false
	})
	if err != nil {
		return nil
	}
	dataformat.SetResToDefault(res, catchErrors) //将template中未找到的未初始化的值根据默认值初始化
	return res
	//logger.Infoln("开始进行模板匹配")
	//resMap := make(map[string]interface{})
	//j.Unmarshal([]byte(template), &resMap)
	////将template封装到map中
	//TravelJsonData([]byte(template), func(path string, value interface{}) bool {
	//	//遍历template
	//	if strValue, ok := value.(string); ok {
	//		if strings.HasPrefix(strValue, "$") { //如果value是以$开头的
	//			if HasIndexMatcher.MatchString(path) { //如果对应的path中有[n] 即存在数组
	//				//Path=$.data.baseInfo[0].detail.personalMonthDepositRatio
	//				//Value=$.result.employeePaymentRate.string
	//				itemCount := getItemCountInJson(body, strValue)
	//				//遍历body 根据路径$.result.employeePaymentRate.string找到body中对应的值出现的个数
	//				if itemCount > 1 { //如果该路径对应的值出现了多次
	//					strIndex := HasIndexMatcher.FindAllStringIndex(path, -1) //[[7 15]] 查看出现数组的是哪一项
	//					splitIndex := strIndex[len(strIndex)-1][1]               //15 取数组中的最后一个数组的第二个值 取最后一个出现数组的那一项
	//					AppendItemInRes(resMap, strings.Split(path[2:splitIndex-3], "."), itemCount)
	//					//传入重复出现的次数 [data baseInfo]和resMap(template) 将最后一个出现数组的那一项包括前面的项传入
	//					//根据重复的次数复制出现数组的那一项并添加到resMap中
	//				}
	//			}
	//		}
	//	}
	//	return false
	//})
	//resData, err := j.Marshal(resMap)
	////>>>>>>> Stashed changes
	//if err != nil {
	//	return nil
	//}
	//logger.Debugf("匹配后的template为:%v", string(resData))
	//resultMap, err := getJsonResultMap(body) //resultMap[path] = value
	////遍历xmlData,将xmlData的path和value存入resultMap中
	//logger.Infof("开始向模板填充返回值")
	//result, err := ChangeValueTravel(resData, func(path string, oldValue interface{}) (interface{}, bool) {
	//	//遍历resData 处理后的template
	//	//log.Println(oldValue)
	//	realValuePath, ok := oldValue.(string)
	//	if !ok {
	//		logger.Errorf("%v格式匹配失败", path)
	//		return nil, false
	//	}
	//	if strings.HasPrefix(realValuePath, "$.") {
	//		pathArray := strings.Split(realValuePath, ".")
	//
	//		valuePath := strings.Join(pathArray[0:len(pathArray)-1], ".") //将value最后的数据类型干掉
	//		var newValue interface{}
	//		//if resultMap == nil {
	//		//	newValue, err = mixformat.go(body, valuePath, pathArray[len(pathArray)-1])
	//		//	if err != nil {
	//		//		return nil, false
	//		//	}
	//		//} else {
	//		//	var ok bool
	//		//	newValue, ok = resultMap[valuePath]
	//		//	if !ok {
	//		//		if IndexMatcher.MatchString(valuePath) {
	//		//			valuePath = IndexMatcher.ReplaceAllString(valuePath, "[0]")
	//		//			newValue, ok = resultMap[valuePath]
	//		//			if !ok {
	//		//				return "", true
	//		//			}
	//		//		} else {
	//		//			return "", true
	//		//		}
	//		//	}
	//		//}
	//		newValue, ok = resultMap[valuePath]
	//		if !ok {
	//			logger.Errorf("返回值中不存在匹配值,匹配值:%v", valuePath)
	//			return nil, false
	//		}
	//		log.Println(newValue)
	//		newTypeValue, err := convertValueType(newValue, pathArray[len(pathArray)-1])
	//		if err == nil {
	//			return newTypeValue, true
	//		}
	//		logger.Errorf("返回值格式转换失败: %s\n", err.Error())
	//		return newValue, true
	//	}
	//	return nil, false
	//})
	//
	//if err != nil {
	//	return nil
	//}
	//
	//logger.Debug(string(result))
	//res := make(map[string]interface{})
	//ioutil.WriteFile("./temp.json", resData, os.ModePerm)
	//err = j.Unmarshal(result, &res)
	//if err != nil {
	//	logger.Error(err)
	//}
	//dataformat.SetResToDefault(res)
	//result, _ = j.Marshal(res)
	//return result
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func GetCommonPrefix(s1 string, s2 string) string {
	//minLen := math.Min(len(s1), len(s2))

	minLen := Min(len(s1), len(s2))
	//res := ""
	index := 0
	for i := 0; i < minLen && s1[i] == s2[i]; i++ {
		//res = append(res, s1[i])
		index = i
	}
	if index > 0 {
		return s1[0 : index+1]
	} else {
		return ""
	}
}


func TransJsonFormatNew(body []byte, template string, args map[string]string, catchError trace.CatchErrors) []byte {
	logger.Infoln("开始进行模板匹配")
	resMap := make(map[string]interface{})
	//ioutil.WriteFile("./template", []byte(template), os.ModePerm)
	err := json.Unmarshal([]byte(template), &resMap)
	if err != nil {
		logger.Errorf("解析template失败:%v", err)
	}
	//将template封装到map中
	pathCount := struct {
		path     string
		strValue string
		count    int
	}{}
	_ = TravelJsonData([]byte(template), func(path string, value interface{}) bool {
		//遍历template
		if strValue, ok := value.(string); ok {
			if strings.HasPrefix(strValue, "$") { //如果value是以$开头的
				if HasIndexMatcher.MatchString(path) { //如果对应的path中有[n] 即存在数组
					//Path=$.data.baseInfo[0].detail.personalMonthDepositRatio
					//Value=$.result.employeePaymentRate.string
					itemCount := getItemCountInJson(body, strValue)
					//遍历body 根据路径$.result.employeePaymentRate.string找到body中对应的值出现的个数
					if itemCount > 1 { //如果该路径对应的值出现了多次
						pathIndex := IndexMatcher.FindAllStringIndex(path, -1) //[[7 15]] 查看出现数组的是哪一项
						splitIndex := pathIndex[len(pathIndex)-1][0]           //15 取数组中的最后一个数组的第二个值 取最后一个出现数组的那一项
						AppendItemInRes(resMap, strings.Split(path[2:splitIndex], "."), itemCount, strValue)
						//传入重复出现的次数 [data baseInfo]和resMap(template) 将最后一个出现数组的那一项包括前面的项传入
						//根据重复的次数复制出现数组的那一项并添加到resMap中
						if len(pathIndex) > 1 {
							for i := 0; i < (len(pathIndex) - 1); i++ {
								temPath := path[:pathIndex[i][0]]
								if pathCount.path == temPath {
									for i := 1; i <= pathCount.count; i++ {
										if strings.HasPrefix(strValue, pathCount.strValue) {

											newPath := path[:len(temPath)] + "[" + strconv.Itoa(i) + "]" + path[len(temPath)+3:]
											newStrValue := strValue[:len(pathCount.strValue)] + "[" + strconv.Itoa(i) + "]" + strValue[len(pathCount.strValue)+3:]
											logger.Debugln(newPath, newStrValue)
											newCount := getItemCountInJson(body, newStrValue)
											if newCount > 1 {
												newPathIndex := IndexMatcher.FindAllStringIndex(newPath, -1) //[[7 15]] 查看出现数组的是哪一项
												newSplitIndex := newPathIndex[len(newPathIndex)-1][0]
												AppendItemInRes(resMap, strings.Split(newPath[2:newSplitIndex], "."), newCount, newStrValue)
											}
										}
									}
								}
							}
							pathCount.path = path[:splitIndex]
							strIndex := IndexMatcher.FindAllStringIndex(strValue, -1)
							pathCount.strValue = strValue[:strIndex[len(strIndex)-1][0]]
							pathCount.count = itemCount
							logger.Debugln(pathCount)
						} else {
							pathCount.path = path[:splitIndex]
							strIndex := IndexMatcher.FindAllStringIndex(strValue, -1)
							pathCount.strValue = strValue[:strIndex[len(strIndex)-1][0]]
							pathCount.count = itemCount
						}
					}
				}
			}
		}
		return false
	})
	resData, err := json.Marshal(resMap)
	//>>>>>>> Stashed changes
	if err != nil {
		return nil
	}
	logger.Debugf("匹配后的template为:%v", string(resData))
	resultMap, err := getJsonResultMap(body) //resultMap[path] = value
	//遍历xmlData,将xmlData的path和value存入resultMap中
	logger.Infof("开始向模板填充返回值")
	result, err := ChangeValueTravel(resData, func(path string, oldValue interface{}) (interface{}, bool) {
		//遍历resData 处理后的template
		//log.Println(oldValue)
		realValuePath, ok := oldValue.(string)
		if !ok {
			return nil, false
		}
		if strings.HasPrefix(realValuePath, "$.") {
			pathArray := strings.Split(realValuePath, ".")
			valuePath := strings.Join(pathArray[0:len(pathArray)-1], ".") //将value最后的数据类型干掉
			logger.Debug(valuePath)
			var newValue interface{}
			newValue, ok = resultMap[valuePath]

			strNewValue, isStr := newValue.(string)
			if isStr && len(strNewValue) <= 0 {
				catchError(trace.FieldIsEmpty, valuePath[2:])
			}
			if !ok {
				if strings.HasPrefix(valuePath, "$.") {
					catchError(trace.FieldDisappeared, valuePath[2:])
				}

				newValue, err = getValueByPath(resultMap, pathArray[1:len(pathArray)-1], pathArray[len(pathArray)-1])
				if err == nil {
					return newValue, true
				}
				logger.Errorf("返回值中不存在匹配值,匹配值:%v", valuePath)
				return nil, false
			}
			newTypeValue, err := convertValueType(newValue, pathArray[len(pathArray)-1])
			if err == nil {
				return newTypeValue, true
			}
			logger.Errorf("返回值格式转换失败: %s\n", err.Error())
			return newValue, false
		} else if strings.HasPrefix(realValuePath, "#.") {
			path := strings.Split(realValuePath, ".")[1]
			if value, ok := args[path]; ok {
				return value, true
			}
			logger.Errorf("参数中不存在匹配值,匹配值:%v", path)
			return nil, false
		}
		return nil, false
	})

	if err != nil {
		return nil
	}

	res := make(map[string]interface{})
	err = json.Unmarshal(result, &res)
	if err != nil {
		logger.Error(err)
	}
	dataformat.SetResToDefault(res, catchError)
	result, _ = json.Marshal(res)
	return result
}

func getItemCountInJson(body []byte, strValue string) int {
	count := 0
	//tempValuePath := IndexMatcher.ReplaceAllString(strValue, "") //$.result.employeePaymentRate.string
	//Path=$.data.baseInfo[0].detail.personalMonthDepositRatio
	//Value=$.result.employeePaymentRate.string
	strIndexArray := IndexMatcher.FindAllStringIndex(strValue, -1)
	var tempValuePath, ralph string
	logger.Debugln(strValue)
	//var isArrayWithArray bool
	if len(strIndexArray) > 0 {
		lastIndex := strIndexArray[len(strIndexArray)-1]
		tempValuePath = strValue[:lastIndex[0]] + "[0]" + strValue[lastIndex[1]:]
		path := strings.Split(tempValuePath, ".")
		ralph = strings.Join(path[0:len(path)-1], ".")
	} else {
		path := strings.Split(strValue, ".")
		ralph = strings.Join(path[0:len(path)-1], ".")
	}
	err := TravelJsonData(body, func(path string, value interface{}) bool {
		pathIndexArray := IndexMatcher.FindAllStringIndex(path, -1)
		var tempPath string
		if len(pathIndexArray) > 0 {
			lastPathIndex := pathIndexArray[len(pathIndexArray)-1]
			tempPath = path[:lastPathIndex[0]] + "[0]" + path[lastPathIndex[1]:]
			//tempPath := IndexMatcher.ReplaceAllString(path, "") //$.result.employeePaymentRate
		} else {
			tempPath = path
		}
		if ralph == tempPath {
			count++
		}
		return false
	})
	if err != nil {
		logger.Error(err)
		return -1
	}
	return count

}
func getJsonResultMap(body []byte) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	//err := json.Unmarshal(body, &res)
	//if err != nil {
	//	return nil, err
	//}

	//err := xml.TraverseXmlNew(xmlData, func(name string, path string, value string) bool {
	//	res[path] = value
	//	return false
	//})
	err := TravelJsonData(body, func(path string, value interface{}) bool {

		res[path] = value
		return false
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//func traverseSourceElement(root map[string]interface{}, source map[string]interface{}, parentPath string, template map[string]interface{}, res map[string]interface{}, arrayCountMap map[string]int) {
//	for k, v := range source {
//		switch value := v.(type) {
//		case map[string]interface{}:
//			traverseSourceElement(root, value, parentPath+"."+k, template, res, arrayCountMap)
//			break
//		case []interface{}:
//			traverseSourceArrayElement(root, value, parentPath+"."+k, template, res, arrayCountMap)
//			break
//		default:
//			//log.Printf("parentPath:%s key:%s\n", parentPath, k)
//			realPath := parentPath + "." + k
//			tempPath := realPath
//			if hasIndexMatcher.MatchString(realPath) {
//				tempPath = string(indexMatcher.ReplaceAllString(realPath, "[0]"))
//			}
//			log.Println(realPath)
//			//setTemplateValue(root, realPath, tempPath, template, v, nil, res, arrayCountMap)
//
//			break
//		}
//	}
//}
//
//func traverseSourceArrayElement(root map[string]interface{}, sourceArray []interface{}, parentPath string, template map[string]interface{}, res map[string]interface{}, arrayCountMap map[string]int) {
//	for index, v := range sourceArray {
//		switch value := v.(type) {
//		case map[string]interface{}:
//			traverseSourceElement(root, value, log.Sprintf("%s[%d]", parentPath, index), template, res, arrayCountMap)
//			break
//		case []interface{}:
//			traverseSourceArrayElement(root, value, log.Sprintf("%s[%d].", parentPath, index), template, res, arrayCountMap)
//			break
//		default:
//			realPath := log.Sprintf("[%d]", index)
//			tempPath := realPath
//			if hasIndexMatcher.MatchString(realPath) {
//				tempPath = string(indexMatcher.ReplaceAllString(realPath, "[0]"))
//			}
//			//setTemplateValue(tempPath,)
//			//setTemplateArrayValue(tempPath,value,v,)
//			//setTemplateValue(root, realPath, tempPath, template, v, nil, res, arrayCountMap)
//			break
//		}
//	}
//}

func setTemplateValue(root map[string]interface{}, realPath string, tempPath string, template map[string]interface{}, realValue interface{}, parent interface{}, res map[string]interface{}) interface{} {
	for k, v := range template {
		switch templateValue := v.(type) {
		case string:
			if strings.HasPrefix(templateValue, "$") {
				if strings.HasPrefix(templateValue, tempPath) && templateValue[len(tempPath)] == '.' {
					parentArray, isArrayItem := parent.([]interface{}) //父亲节点是否是jsonArray

					//log.Println(isArrayItem)
					if isArrayItem {
						isNeedNewItem := true
						for parentItemIndex, resArrayItem := range parentArray {
							item := resArrayItem.(map[string]interface{})
							if itemValue, ok := item[k]; ok {
								if stringValue, success := itemValue.(string); success {
									if strings.HasPrefix(stringValue, "$") {
										if strings.HasSuffix(stringValue, "e%x") {
											isNeedNewItem = false
											//res[k] = realPath
											//continue
											if len(parentArray) == 1 {
												path := strings.Split(templateValue, ".")
												res[k], _ = convertValueType(realValue, path[len(path)-2])
											}
											continue
										}
										//log.Printf("-----------------%s\n", realPath)
										path := strings.Split(templateValue, ".")
										itemRealValue, err := convertValueType(realValue, path[len(path)-1])
										if err == nil {
											item[k] = itemRealValue
										}
										//item[k], err := convertValueType(realValue, path[len(path)-1])

										//log.Printf("key : %s,value : %v", k, item[k])
										parentArray[parentItemIndex] = item
										isNeedNewItem = false

										for resItemKey, resItemValue := range res {
											if itemStringValue, success := resItemValue.(string); success {
												if strings.HasPrefix(itemStringValue, "$") && strings.HasSuffix(itemStringValue, "e%x") {
													commonPrefix := GetCommonPrefix(tempPath, itemStringValue)
													if len(commonPrefix) > 0 {
														path := strings.Split(itemStringValue, ".")
														exlen := 0
														if strings.HasSuffix(itemStringValue, "e%x") {
															path = path[0 : len(path)-1]
															exlen = 4
														}
														newRealPath := realPath[0:len(commonPrefix)] + itemStringValue[len(commonPrefix):len(itemStringValue)-exlen]
														emptyValue, err := getValueByPath(root, strings.Split(newRealPath, ".")[1:len(path)-1], path[len(path)-1])
														//log.Printf("realPath: %s newRealPath: %s, itemStringValue:%s,err:%s", realPath, newRealPath, itemStringValue, err)
														if err != nil {
															res[resItemKey] = nil
														} else {
															res[resItemKey] = emptyValue
														}
													} else {
														res[resItemKey] = nil
													}
												}
											}
										}
										break
									}
								}
							}
						}

						if isNeedNewItem {
							newItem := copyJsonItem(template)
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
						path := strings.Split(templateValue, ".")
						res[k], _ = convertValueType(realValue, path[len(path)-1])
						continue
					}
				}
			}
			break
		case []interface{}:
			//length := len(templateValue)
			resArray := res[k].([]interface{})
			setTemplateArrayValue(root, realPath, tempPath, templateValue, realValue, k, res, resArray)
			break
		case map[string]interface{}:
			resMap := res[k].(map[string]interface{})
			setTemplateValue(root, realPath, tempPath, templateValue, realValue, res, resMap)
			break
		}
	}
	return nil
}

func getValueByPath(source map[string]interface{}, path []string, valueType string) (interface{}, error) {
	for index, key := range path {
		isArray := false
		var tempKey string
		var indexStr string
		if IndexMatcher.MatchString(key) {
			tempKey = IndexMatcher.ReplaceAllString(key, "")
			indexInKey := IndexMatcher.FindString(key)
			indexStr = indexInKey[1 : len(indexInKey)-1]
			isArray = true
		} else {
			tempKey = key
		}

		switch value := source[tempKey].(type) {
		case map[string]interface{}:
			if isArray {
				return nil, errors.New(fmt.Sprintf("%s is array,but read is jsonobject", key))
			}
			return getValueByPath(value, path[index+1:], valueType)
		case []interface{}:
			if !isArray {
				return nil, errors.New(fmt.Sprintf("%s is not array,but read is array", key))
			}

			//indexStr := key[len(key)-2 : len(key)-1]
			sourceIndex, err := strconv.Atoi(indexStr)
			//log.Printf("sourceIndex:%d,%s\n", sourceIndex, path[index+1:])
			if err != nil {
				return nil, err
			}
			return getValueByPathInArray(value, sourceIndex, path[index+1:], valueType)
		default:
			if isArray {
				return nil, errors.New(fmt.Sprintf("%s is array,but read is Value", key))
			}
			return convertValueType(value, valueType)
		}
	}
	return nil, errors.New("path is empty")
}

func getValueByPathInArray(array []interface{}, index int, path []string, valueType string) (interface{}, error) {
	value := array[index]
	switch v := value.(type) {
	case map[string]interface{}:
		return getValueByPath(v, path, valueType)
	case []interface{}:
		key := path[0]
		if pureIndexMatcher.MatchString(key) {
			index, _ := strconv.Atoi(key[1 : len(key)-1])
			return getValueByPathInArray(v, index, path[1:], valueType)
		} else {
			return nil, errors.New(fmt.Sprintf("key %s must be [num] format", key))
		}
	default:
		logger.Debug(reflect.TypeOf(v))
		return v, nil
	}
}

func setTemplateArrayValue(root map[string]interface{}, realPath string, tempPath string, template []interface{}, realValue interface{}, parentKey string, parent interface{}, res []interface{}) {
	for _, v := range template {
		switch value := v.(type) {
		case map[string]interface{}:
			resItem := res[len(res)-1].(map[string]interface{})
			newArray := setTemplateValue(root, realPath, tempPath, value, realValue, res, resItem)
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
			setTemplateArrayValue(root, realPath, tempPath, value, realValue, "", res, parentValue)
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

func AppendItemInRes(res map[string]interface{}, pathArray []string, count int, itemPath string) error {

	/*
		res   :模板
		pathArray:返回给后台索引字符数组   仅一个值
		itemCount:数组长度
		strValue :从第三方取值的索引字符串   仅一个值
	*/

	//pathArray = [data[0] baseInfo] count=2
	path := pathArray[0] //data[0]
	index := 0
	isArray := false
	if IndexMatcher.MatchString(path) {
		isArray = true
		indexStr := IndexMatcher.FindString(path)
		index, _ = strconv.Atoi(indexStr[1 : len(indexStr)-1])
		path = path[0:strings.Index(path, "[")]
	}
	switch value := res[path].(type) {
	case map[string]interface{}:
		if isArray {
			return errors.New("bad format1")
		} else {
			if len(pathArray) == 1 {
				return errors.New("bad format2")
			} else {
				return AppendItemInRes(value, pathArray[1:], count, itemPath)
			}
		}
	case []interface{}:
		if len(pathArray) == 1 {
			if len(value) < count {
				logger.Debugf("trace:%d", count)
				for c := len(value); c < count; c++ {
					item := copyJsonItemWithIndex(value[0].(map[string]interface{}), c, itemPath)
					value = append(value, item)
				}
				res[path] = value
			}
			return nil
		} else {
			if !isArray {
				return errors.New("bad format3")
			} else {
				return appendItemInArray(value, index, pathArray[1:], count, itemPath)
			}
		}
	default:
		return errors.New("bad format6")
	}
}

func copyJsonItemWithIndex(template map[string]interface{}, index int, commonPath string) map[string]interface{} {
	//return getInsertValues(template,template)
	res := make(map[string]interface{})
	for key, value := range template {
		switch realValue := value.(type) {
		case map[string]interface{}:
			res[key] = copyJsonItemWithIndex(realValue, index, commonPath)
			break
		case []interface{}:
			res[key] = copyJsonArrayItemWithIndex(realValue, index, commonPath)
			break
		case string:
			logger.Debugf("commonPath:%s", commonPath)
			logger.Debugf("realValue:%s", realValue)
			if strings.HasPrefix(realValue, "$") {
				commonPrefix := GetCommonPrefix(commonPath, realValue)
				if HasIndexMatcher.MatchString(commonPrefix) {
					logger.Debug("match1")
					allIndex := IndexMatcher.FindAllStringIndex(commonPrefix, -1)
					lastIndex := allIndex[len(allIndex)-1]
					res[key] = realValue[:lastIndex[0]] + "[" + strconv.Itoa(index) + "]" + realValue[lastIndex[1]:]
				} else if pureIndexMatcher.MatchString(commonPrefix) {
					logger.Debug("match2")
					allIndex := IndexMatcher.FindAllStringIndex(commonPrefix, -1)
					lastIndex := allIndex[len(allIndex)-1]
					res[key] = realValue[:lastIndex[0]] + "[" + strconv.Itoa(index) + "]" + realValue[lastIndex[1]:]
				} else {
					res[key] = realValue
				}
			} else {
				res[key] = realValue
			}
			break
		default:
			res[key] = realValue
			break
		}
	}
	return res
}

func copyJsonArrayItemWithIndex(template []interface{}, index int, commonPath string) []interface{} {
	res := make([]interface{}, 0, 10)
	switch realValue := template[0].(type) {
	case map[string]interface{}:
		res = append(res, copyJsonItemWithIndex(realValue, index, commonPath))
		break
	case []interface{}:
		res = append(res, copyJsonArrayItemWithIndex(realValue, index, commonPath))
		break
	case string:
		if strings.HasPrefix(realValue, "$") {
		} else {
			res = append(res, realValue)
		}
		break
	default:
		res = append(res, realValue)
		break
	}
	return res
}

func appendItemInArray(res []interface{}, index int, pathArray []string, count int, itemPath string) error {
	switch value := res[index].(type) {
	case map[string]interface{}:
		return AppendItemInRes(value, pathArray, count, itemPath)
	case []interface{}:

		//TODO             ???????????????
		return errors.New("bad format4")
	default:
		return errors.New("bad format5")
	}
}

func copyJsonItem(template map[string]interface{}) map[string]interface{} {
	//return getInsertValues(template,template)
	res := make(map[string]interface{})
	for key, value := range template {
		switch realValue := value.(type) {
		case map[string]interface{}:
			res[key] = copyJsonItem(realValue)
			break
		case []interface{}:
			res[key] = copyJsonArrayItem(realValue)
			break
		case string:
			if strings.HasPrefix(realValue, "$") {
				res[key] = realValue
			} else {
				res[key] = realValue
			}
			break
		default:
			res[key] = realValue
			break
		}
	}
	return res
}

func copyJsonArrayItem(template []interface{}) []interface{} {
	res := make([]interface{}, 0, 10)
	if len(template) == 0 {
		return res
	}
	switch realValue := template[0].(type) {
	case map[string]interface{}:
		res = append(res, copyJsonItem(realValue))
		break
	case []interface{}:
		res = append(res, copyJsonArrayItem(realValue))
		break
	case string:
		if strings.HasPrefix(realValue, "$") {
		} else {
			res = append(res, realValue)
		}
		break
	default:
		res = append(res, realValue)
		break
	}
	return res
}

func CopyJsonItem(template map[string]interface{}) map[string]interface{} {
	return copyJsonItem(template)
}

func GetValueFromJson(source []byte, needPath string) interface{} {
	var res interface{}
	TravelJsonData(source, func(path string, value interface{}) bool {
		if path == needPath {
			res = value
			return true
		}
		return false
	})
	return res
}

func convertValueType(value interface{}, valueType string) (interface{}, error) {
	if valueType == "string" {
		if value == nil {
			return "", nil
		}
		switch v := value.(type) {

		case float64:
			return fmt.Sprintf("%f", v), nil
		default:
			return fmt.Sprintf("%v", value), nil
		}
	}
	switch v := value.(type) {
	case float64:
		if valueType == "float" {
			return v, nil
		} else if valueType == "int" {
			return int(v), nil
		} else if valueType == "string" {
			return fmt.Sprintf("%f", v), nil
		} else {
			logger.Errorln("格式转换失败")
			return nil, errors.New(fmt.Sprintf("%s type not favor", valueType))
		}
	case string:
		if valueType == "float" {
			return strconv.ParseFloat(v, 64)
		} else if valueType == "int" {
			return strconv.ParseInt(v, 10, 64)
		} else if valueType == "string" {
			return v, nil
		} else {
			logger.Errorln("格式转换失败")
			return nil, errors.New(fmt.Sprintf("%s type not favor", valueType))
		}
	default:
		if v == nil {
			switch valueType {
			case "float":
				return 0, nil
			case "int":
				return 0, nil
			case "string":
				return "", nil
			}
		}
		return v, nil
	}
}

func IsJsonMathCondition(source []byte, config []*dataformat.MatchConfig) bool {
	ExistMap := make(map[string]bool)
	isMatch := true

	for index, value := range config {
		if len(source) == 0 && value.Condition == "isnil" {
			logger.Debugf("返回值匹配成功:%v", config[index])
			config = append(config[:index], config[index+1:]...)
			return isMatch
		}
	}
	err := TravelJsonData(source, func(path string, value interface{}) bool {
		if len(config) == 0 {
			return true
		}
		for i, c := range config {
			switch c.Condition {
			case "=":
				if path == c.FieldPath {
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
				if path == c.FieldPath {
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
				if strings.Contains(path, c.FieldPath) {
					logger.Debugf("返回值匹配成功:%v", config[i])
					config = append(config[:i], config[i+1:]...)
					return false
				}
				break
			case "!exist":
				if strings.Contains(path, c.FieldPath) {
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

//func findSameInJsonBody(template map[string]interface{}, body map[string]interface{}, res bool) bool {
//	if !res {
//		return res
//	}
//	for key, value := range template {
//		switch value.(type) {
//		case map[string]interface{}:
//			bodyValue, ok := body[key].(map[string]interface{})
//			res = findSameInJsonBody(value.(map[string]interface{}), bodyValue, ok)
//			break
//		case []interface{}:
//			bodyValue, ok := body[key].([]interface{})
//			res = findSameInArrayBody(value.([]interface{}), bodyValue, ok)
//			break
//		default:
//			res = res && reflect.TypeOf(value) == reflect.TypeOf(body[key]) && body[key] == value
//			if !res {
//				return res
//			}
//		}
//	}
//	return res
//}
//
//func findSameInArrayBody(template []interface{}, body []interface{}, res bool) bool {
//	//res := true
//	for index, v := range template {
//		switch v.(type) {
//		case []interface{}:
//			bodyValue, ok := body[index].([]interface{})
//			res = findSameInArrayBody(v.([]interface{}), bodyValue, ok)
//			break
//		case map[string]interface{}:
//			bodyValue, ok := body[index].(map[string]interface{})
//			res = findSameInJsonBody(v.(map[string]interface{}), bodyValue, ok)
//			break
//		default:
//			res = reflect.TypeOf(v) == reflect.TypeOf(body[index]) && body[index] == v
//			if !res {
//				return res
//			}
//		}
//	}
//	return res
//}

//func generateResultByTemplate(template map[string]interface{}, res map[string]interface{}, parent interface{}, source map[string]interface{}) {
//	for k, v := range template {
//		switch value := v.(type) {
//		case map[string]interface{}:
//			resItem := res[k].(map[string]interface{})
//			break
//		case []interface{}:
//			resArray := res[k].([]interface{})
//
//			break
//		case string:
//			if strings.HasPrefix(value, "$") {
//				//reflect.ValueOf(parent)
//				//parentArray, isArray := (*parent).([]interface{})
//				isArray := false
//				pValue := reflect.ValueOf(parent)
//				pValueElement := pValue.Elem()
//				if pValueElement.Kind() == reflect.Slice {
//					isArray = true
//				}
//
//				if isArray {
//					if hasIndexMatcher.MatchString(value) {
//						matchStrIndex := indexMatcher.FindAllStringIndex(value, -1)
//						arrayRealPath := value[2:matchStrIndex[len(matchStrIndex)-1][1]]
//						path := strings.Split(arrayRealPath, ".")
//						count := getArrayCount(path[:len(path)-1], source)
//						if count > -1 {
//							if pValueElement.Len() < count {
//								for n := 1; n < count; n++ {
//									if n >= pValueElement.Cap() {
//										newcap := pValueElement.Cap() + pValueElement.Cap()/2
//										if newcap < 4 {
//											newcap = 4
//										}
//										newv := reflect.MakeSlice(pValueElement.Type(), pValueElement.Len(), newcap)
//										reflect.Copy(newv, v)
//										pValueElement.Set(newv)
//									}
//									if n >= pValueElement.Len() {
//										pValueElement.SetLen(n + 1)
//									}
//									item := pValueElement.Index(n)
//									newItem := copyJsonItemWithIndex(template, "$."+arrayRealPath, n)
//									//parentArray = append(parentArray, item)
//									item.Set(reflect.ValueOf(newItem))
//								}
//							}
//						}
//					}
//				}
//			}
//			break
//		}
//	}
//}
//type jsonType int
//
//const (
//	mapType jsonType = iota
//	arrType
//	iteral
//)
//
//type jsonNode struct {
//	parent *jsonNode
//	value  interface{}
//	key    interface{}
//	length int
//}
//
//func (jNode *jsonNode) itemType() jsonType {
//	switch jNode.value.(type) {
//	case map[string]interface{}:
//		return mapType
//	case []interface{}:
//		return arrType
//	default:
//		return iteral
//	}
//}

//func generateResultByTemplate(template map[string]interface{}, res map[string]interface{}, node *jsonNode, source map[string]interface{}) {
//	for k, v := range template {
//		switch value := v.(type) {
//		case map[string]interface{}:
//			resItem := res[k].(map[string]interface{})
//			node := jsonNode{node, resItem, k, -1}
//			generateResultByTemplate(value, resItem, &node, source)
//		case []interface{}:
//			resArray := res[k].([]interface{})
//			node := jsonNode{node, resArray, k, -1}
//			generateArrayByTemplate(value, resArray, &node, source)
//		case string:
//			if strings.HasPrefix(value, "$") {
//				//parentArray, isArray := (*parent).([]interface{})
//				//if node.itemType() == arrType {
//				//
//				//}
//				if node.itemType() == arrType && hasIndexMatcher.MatchString(value) {
//					parentArray := node.value.([]interface{})
//					if hasIndexMatcher.MatchString(value) {
//						matchStrIndex := indexMatcher.FindAllStringIndex(value, -1)
//						arrayRealPath := value[2:matchStrIndex[len(matchStrIndex)-1][1]]
//						path := strings.Split(arrayRealPath, ".")
//						count := getArrayCount(path[:len(path)-1], source)
//						if count > -1 {
//							if len(parentArray) == 1 {
//								for n := 1; n < count; n++ {
//									newItem := copyJsonItemWithIndex(template, "$."+arrayRealPath, n, source)
//									//generateResultByTemplate(copyJsonItem(newItem), newItem, &jsonNode{nil, newItem, nil}, source)
//									parentArray = append(parentArray, newItem)
//								}
//								switch key := node.key.(type) {
//								case int:
//									parentParentArray := node.parent.value.([]interface{})
//									parentParentArray[key] = parentArray
//									break
//								case string:
//									parentMap := node.parent.value.(map[string]interface{})
//									parentMap[key] = parentArray
//									break
//								default:
//									break
//								}
//							} else {
//
//							}
//						}
//					}
//				}
//
//			}
//			break
//		}
//	}
//}
//
//func generateArrayByTemplate(template []interface{}, res []interface{}, node *jsonNode, source map[string]interface{}) {
//	for index, v := range template {
//		switch value := v.(type) {
//		case map[string]interface{}:
//			resItem := res[index].(map[string]interface{})
//			generateResultByTemplate(value, resItem, &jsonNode{node, resItem, index, -1}, source)
//		case []interface{}:
//			resArray := res[index].([]interface{})
//			generateArrayByTemplate(value, resArray, &jsonNode{node, resArray, index, -1}, source)
//		case string:
//
//			break
//		}
//	}
//}

//func copyJsonItemWithIndex(template map[string]interface{}, realPath string, index int, source map[string]interface{}) (map[string]interface{}) {
//	//return getInsertValues(template,template)
//	isNeedGenerate := false
//	res := make(map[string]interface{})
//	for key, value := range template {
//		switch valueInTemplate := value.(type) {
//		case map[string]interface{}:
//			res[key] = copyJsonItemWithIndex(valueInTemplate, realPath, index, source)
//			break
//		case []interface{}:
//			res[key] = copyJsonArrayItemWithIndex(valueInTemplate, realPath, index, source)
//			break
//		case string:
//
//			if strings.HasPrefix(valueInTemplate, "$") {
//				if hasIndexMatcher.MatchString(valueInTemplate) {
//					//hasIndexMatcher.ReplaceAllStringFunc(realValue, func(s string) string {
//					//
//					//})
//					tempPath := hasIndexMatcher.ReplaceAllString(realPath, "[0]")
//					if strings.HasPrefix(valueInTemplate, tempPath) {
//						prefix := GetCommonPrefix(valueInTemplate, tempPath)
//						leftIndex := strings.LastIndex(realPath, "[")
//						newValue := realPath[:leftIndex] + "[" + strconv.Itoa(index) + "]" + valueInTemplate[len(prefix):]
//						if hasIndexMatcher.MatchString(valueInTemplate[len(prefix):]) {
//							isNeedGenerate = true
//						}
//						res[key] = newValue
//					} else {
//						res[key] = valueInTemplate
//					}
//				} else {
//					res[key] = valueInTemplate
//				}
//			} else {
//				res[key] = valueInTemplate
//			}
//			break
//		default:
//			res[key] = valueInTemplate
//			break
//		}
//	}
//
//	if isNeedGenerate {
//		generateResultByTemplate(copyJsonItem(res), res, &jsonNode{nil, res, nil, -1}, source)
//	}
//
//	return res
//}
//
//func copyJsonArrayItemWithIndex(template []interface{}, realPath string, index int, source map[string]interface{}) []interface{} {
//	res := make([]interface{}, 0, 10)
//
//	switch templateValue := template[0].(type) {
//	case map[string]interface{}:
//		res = append(res, copyJsonItemWithIndex(templateValue, realPath, index, source))
//		break
//	case []interface{}:
//		res = append(res, copyJsonArrayItemWithIndex(templateValue, realPath, index, source))
//		break
//	case string:
//		if strings.HasPrefix(templateValue, "$") {
//			if hasIndexMatcher.MatchString(templateValue) {
//				tempPath := hasIndexMatcher.ReplaceAllString(realPath, "[0]")
//				if strings.HasPrefix(templateValue, tempPath) {
//					prefix := GetCommonPrefix(templateValue, tempPath)
//					leftIndex := strings.LastIndex(realPath, "[")
//					newValue := realPath[:leftIndex] + "[" + strconv.Itoa(index) + "]" + templateValue[len(prefix):]
//
//					//res[key] = newValue
//					res = append(res, newValue)
//				} else {
//					res = append(res, templateValue)
//				}
//			} else {
//				res = append(res, templateValue)
//			}
//		} else {
//			res = append(res, templateValue)
//		}
//		break
//	default:
//		res = append(res, templateValue)
//		break
//	}
//	return res
//}

//func getArrayCountNew(path []string, node *jsonNode, realPath string, source map[string]interface{}) int {
//	if hasIndexMatcher.MatchString(realPath) {
//		arrayIndex := hasIndexMatcher.FindAllString(realPath, -1)
//		for n := len(arrayIndex) - 1; n >= 0; n-- {
//			arrayItem := arrayIndex[n]
//
//			indexInPath := getIndexInPath(arrayItem, path)
//			if indexInPath >= 0 {
//
//			}
//		}
//	}
//}

//func getIndexInPath(itemPath string, path []string) int {
//	for i, p := path {
//		if itemPath == p {
//			return i
//		}
//	}
//	return -1
//}

//func getArrayCount(path []string, source map[string]interface{}) int {
//	for index, key := range path {
//		isArray := false
//		var tempKey string
//		var indexStr string
//		if indexMatcher.MatchString(key) {
//			tempKey = indexMatcher.ReplaceAllString(key, "")
//			indexInKey := indexMatcher.FindString(key)
//			indexStr = indexInKey[1 : len(indexInKey)-1]
//			isArray = true
//		} else {
//			tempKey = key
//		}
//
//		switch value := source[tempKey].(type) {
//		case map[string]interface{}:
//			if isArray {
//				return -1
//			}
//			return getArrayCount(path[index+1:], value)
//		case []interface{}:
//			if !isArray {
//				return -1
//			}
//			if index == len(path)-1 {
//				return len(value)
//			} else {
//				//indexStr := key[len(key)-2 : len(key)-1]
//				sourceIndex, err := strconv.Atoi(indexStr)
//				if err != nil {
//					return -1
//				}
//				return getArrayCountInArray(value, sourceIndex, path[index+1:])
//			}
//		default:
//
//			log.Println(reflect.TypeOf(value), tempKey)
//			break
//		}
//	}
//	return -1
//}
//
//func getArrayCountInArray(array []interface{}, index int, path []string) int {
//	v := array[index]
//	switch value := v.(type) {
//	case map[string]interface{}:
//		return getArrayCount(path, value)
//		//break
//	case []interface{}:
//		key := path[0]
//		if len(path) == 1 {
//			return len(value)
//		} else {
//			if pureIndexMatcher.MatchString(key) {
//				index, _ := strconv.Atoi(key[1 : len(key)-1])
//				return getArrayCountInArray(value, index, path[1:])
//			} else {
//				return -1
//			}
//		}
//	}
//	return -1
//}
