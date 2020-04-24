package dataformat

import (
	"github.com/yanzongzhen/DataFormatUtils/trace"
	"strings"
)

func SetResToDefault(res map[string]interface{}, catchErrors trace.CatchErrors) {
	for k, v := range res {
		switch value := v.(type) {
		case []interface{}:
			SetResArrayToDefault(value, catchErrors)
			break
		case map[string]interface{}:
			SetResToDefault(value, catchErrors)
			break
		case string:
			if strings.HasPrefix(value, "$") {

				catchErrors(trace.FieldDisappeared, value[2:])
				valueType := ""
				path := strings.Split(value, ".")
				if strings.HasSuffix(value, "e%x") {
					valueType = path[len(path)-2]
				} else {
					valueType = path[len(path)-1]
				}
				switch valueType {
				case "string":
					res[k] = ""
					break
				case "float":
					res[k] = 0.0
					break
				case "int":
					res[k] = 0
					break
				default:
					res[k] = nil
					break
				}
			}
			break
		default:
			break
		}
	}
}

func SetResArrayToDefault(res []interface{}, catchErrors trace.CatchErrors) {
	for index, v := range res {
		switch value := v.(type) {
		case []interface{}:
			SetResArrayToDefault(value, catchErrors)
			break

		case map[string]interface{}:
			SetResToDefault(value, catchErrors)
			break
		case string:
			if strings.HasPrefix(value, "$") {
				valueType := ""
				path := strings.Split(value, ".")
				if strings.HasSuffix(value, "e%x") {
					valueType = path[len(path)-2]
				} else {
					valueType = path[len(path)-1]
				}
				switch valueType {
				case "string":
					res[index] = ""
					break
				case "float":
					res[index] = 0.0
					break
				case "int":
					res[index] = 0
					break
				default:
					res[index] = nil
					break
				}
			}
			break
		default:
			break
		}
	}
}
