package traz

import (
	"gopkg.in/yaml.v2"
)

type YAMLtraz struct {
	data string
}

func NewYAMLtraz(data string) *YAMLtraz {
	return &YAMLtraz{
		data: data,
	}
}

func (tt *YAMLtraz) Apply(s interface{}) error {
	return yaml.Unmarshal([]byte(tt.data), s)
}
