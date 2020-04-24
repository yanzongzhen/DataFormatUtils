package mixformat

type GetValue func(path []string, dataType string) (interface{}, error)

//func JsonToJson(jsonData []byte, template string) ([]byte, error) {
//	templateMap := make(map[string]interface{})
//	err := j.Unmarshal([]byte(template), &templateMap)
//	if err != nil {
//		return nil, err
//	}
//	res := json.CopyJsonItem(templateMap)
//	err = json.TravelJsonData(jsonData, func(path string, value interface{}) {
//
//		setTemplateValue(jsonData, path, path, templateMap, value, nil, res)
//	})
//
//	if err != nil {
//		return nil, err
//	}
//	setResToDefault(res)
//	return j.Marshal(res)
//}
