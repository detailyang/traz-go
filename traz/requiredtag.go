package traz

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
)

const (
	DefaultRequiredTagName = "required"
)

type RequiredTagtraz struct {
	tagName string
}

func NewRequiredTagtraz(tagName string) *RequiredTagtraz {
	return &RequiredTagtraz{
		tagName: tagName,
	}
}

func (rt *RequiredTagtraz) Apply(s interface{}) error {
	for _, field := range structs.Fields(s) {
		if err := rt.apply("", field); err != nil {
			return err
		}
	}
	return nil
}

func (rt *RequiredTagtraz) apply(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		fieldName += "."

		for _, f := range field.Fields() {
			if err := rt.apply(fieldName, f); err != nil {
				return err
			}
		}

	default:
		val := field.Tag(rt.tagName)
		if val == "" {
			return nil
		}

		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("traz: invalid value %q", val)
		}

		if !b { // not required
			return nil
		}

		if field.IsZero() {
			return fmt.Errorf("traz: field '%s' is required", fieldName)
		}
	}

	return nil
}
