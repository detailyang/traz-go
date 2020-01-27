package traz

import (
	"github.com/BurntSushi/toml"
)

type TOMLtraz struct {
	data string
}

func NewTOMLtraz(data string) *TOMLtraz {
	return &TOMLtraz{
		data: data,
	}
}

func (tt *TOMLtraz) Apply(s interface{}) error {
	_, err := toml.Decode(tt.data, s)
	return err
}
