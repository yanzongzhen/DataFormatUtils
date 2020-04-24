package dataformat

type MatchConfig struct {
	FieldPath string      `json:"fields_path"`
	Condition string      `json:"condition"`
	Value     interface{} `json:"value"`
}
