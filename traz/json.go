package traz

type JSONtraz = YAMLtraz

func NewJSONtraz(data string) *JSONtraz {
	return NewYAMLtraz(data)
}
