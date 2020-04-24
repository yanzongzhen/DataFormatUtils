/**
 * @Author: guanyunlong
 * @Description:
 * @File:  trace
 * @Version: 1.0.0
 * @Date: 20-3-4 下午2:26
 */
package trace

var (
	FieldDisappeared = "fieldDisappeared"
	FieldIsEmpty     = "fieldIsEmpty"
)

type CatchErrors func(errorType string, fieldName string)
