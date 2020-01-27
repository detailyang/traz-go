package traz

import (
	"reflect"

	"github.com/fatih/structs"
)

const (
	DefaultDefaultTagName = "default"
)

type DefaultTagtraz struct {
	tagName string
}

func NewDefaultTagtraz(tagName string) *DefaultTagtraz {
	return &DefaultTagtraz{
		tagName: tagName,
	}
}

func (sv *DefaultTagtraz) Apply(s interface{}) error {
	if sv.tagName == "" {
		sv.tagName = DefaultDefaultTagName
	}

	for _, field := range structs.Fields(s) {
		if err := sv.apply(sv.tagName, field); err != nil {
			return err
		}
	}

	return nil
}

func (t *DefaultTagtraz) apply(tagName string, field *structs.Field) error {
	switch field.Kind() {
	case reflect.Struct:
		switch field.Value().(type) {
		case Duration, Size:
			return t.set(t.tagName, field)
		}

		for _, f := range field.Fields() {
			if err := t.apply(tagName, f); err != nil {
				return err
			}
		}
	default:
		return t.set(t.tagName, field)
	}

	return nil
}

func (t *DefaultTagtraz) set(tagName string, field *structs.Field) error {
	defaultVal := field.Tag(t.tagName)
	if defaultVal == "" {
		return nil
	}

	return fieldSet(field, defaultVal)
}
