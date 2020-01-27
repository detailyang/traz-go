package traz

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type StructConstructor func(prefix string, name string, m interface{}) (interface{}, bool, error)

type MapPasstraz struct {
	constructors []StructConstructor
}

func NewMapPasstraz(constructors ...StructConstructor) *MapPasstraz {
	return &MapPasstraz{
		constructors: constructors,
	}
}

func (mt *MapPasstraz) Append(constructors ...StructConstructor) {
	mt.constructors = append(mt.constructors, constructors...)
}

func (mt *MapPasstraz) Apply(s interface{}) error {
	prefix := mt.genprefix("", structs.Name(s))
	for _, field := range structs.Fields(s) {
		if err := mt.apply(prefix, field); err != nil {
			return err
		}
	}
	return nil
}

func (mt *MapPasstraz) pass(prefix, name string, input interface{}) (interface{}, bool, error) {
	for i := range mt.constructors {
		o, ok, err := mt.constructors[i](prefix, name, input)
		if err != nil {
			return nil, ok, err
		}
		if !ok {
			continue
		}
		return o, true, nil
	}
	return nil, false, nil
}

func (mt *MapPasstraz) apply(prefix string, field *structs.Field) error {
	prefix = mt.genprefix(prefix, field.Name())
	switch f := field.Value().(type) {
	case map[string]interface{}:
		for k, v := range f {
			o, ok, err := mt.pass(mt.genprefix(prefix, k), k, v)
			if err != nil {
				return err
			}
			if ok {
				f[k] = o
			}
		}
	default:
		switch field.Kind() {
		case reflect.Struct:
			for _, f := range field.Fields() {
				if err := mt.apply(prefix, f); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (mt *MapPasstraz) genprefix(prefix, name string) string {
	if prefix == "" {
		return strings.ToUpper(name)
	}
	return strings.ToUpper(prefix) + "." + strings.ToUpper(name)
}
