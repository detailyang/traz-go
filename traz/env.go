package traz

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type EnvGetter interface {
	Getenv(key string) string
}

type EnvGetterFunc func(key string) string

func (e EnvGetterFunc) Getenv(key string) string {
	return e(key)
}

type Envtraz struct {
	prefix string
	getter EnvGetter
}

func NewEnvtraz(prefix string) *Envtraz {
	return &Envtraz{
		prefix: prefix,
		getter: EnvGetterFunc(os.Getenv),
	}
}

func (et *Envtraz) Apply(s interface{}) error {
	prefix := et.genprefix(et.prefix, structs.Name(s))
	for _, field := range structs.Fields(s) {
		if err := et.apply(prefix, field); err != nil {
			return err
		}
	}
	return nil
}

func (et *Envtraz) String(w io.Writer, s interface{}) error {
	prefix := et.genprefix(et.prefix, structs.Name(s))
	for _, field := range structs.Fields(s) {
		if err := et.string(w, prefix, field); err != nil {
			return err
		}
	}
	return nil
}

func (et *Envtraz) string(w io.Writer, prefix string, field *structs.Field) error {
	prefix = et.genprefix(prefix, field.Name())

	switch field.Kind() {
	case reflect.Struct:
		switch field.Value().(type) {
		case Duration, Size:
			return et.fprintf(w, prefix, field)
		}

		for _, f := range field.Fields() {
			if err := et.string(w, prefix, f); err != nil {
				return err
			}
		}
	default:
		switch f := field.Value().(type) {
		case map[string]interface{}:
			return et.stringMapStringInterface(w, prefix, f)
		}
		if err := et.fprintf(w, prefix, field); err != nil {
			return err
		}
	}

	return nil
}

func (et *Envtraz) fprintf(w io.Writer, prefix string, field *structs.Field) error {
	kind := field.Kind().String()
	switch field.Value().(type) {
	case Duration:
		kind = "duration"
	case Size:
		kind = "size"
	}
	_, err := fmt.Fprintf(w, "      %s %s Change value of %q.\n",
		prefix, kind, prefix)
	return err
}

func (et *Envtraz) apply(prefix string, field *structs.Field) error {
	prefix = et.genprefix(prefix, field.Name())
	switch field.Kind() {
	case reflect.Struct:
		switch field.Value().(type) {
		case Duration, Size:
			return et.set(prefix, field)
		}

		for _, f := range field.Fields() {
			if err := et.apply(prefix, f); err != nil {
				return err
			}
		}
	default:
		switch f := field.Value().(type) {
		case map[string]interface{}:
			return et.applyMapStringInterface(prefix, f)
		}
		if err := et.set(prefix, field); err != nil {
			return err
		}
	}

	return nil
}

func (et *Envtraz) stringMapStringInterface(w io.Writer, prefix string, m map[string]interface{}) error {
	for key := range m {
		if structs.IsStruct(m[key]) {
			if structs.IsStruct(m[key]) {
				mprefix := et.genprefix(prefix, key)
				for _, field := range structs.Fields(m[key]) {
					mmprefix := et.genprefix(mprefix, field.Name())
					if err := et.fprintf(w, mmprefix, field); err != nil {
						return err
					}
				}
			}
			continue

		} else {
			// reflect.
		}
	}
	return nil
}

func (et *Envtraz) applyMapStringInterface(prefix string, m map[string]interface{}) error {
	for key := range m {
		if structs.IsStruct(m[key]) {
			mprefix := et.genprefix(prefix, key)
			if structs.IsStruct(m[key]) {
				for _, field := range structs.Fields(m[key]) {
					if err := et.apply(mprefix, field); err != nil {
						return err
					}
				}
			}
			continue

		} else {
			// reflect.
		}
	}
	return nil
}

func (et *Envtraz) set(prefix string, field *structs.Field) error {
	v := et.getter.Getenv(prefix)
	if v == "" {
		return nil
	}

	return fieldSet(field, v)
}

func (et *Envtraz) genprefix(prefix string, name string) string {
	if prefix == "" {
		return strings.ToUpper(name)
	}
	return strings.ToUpper(prefix) + "." + strings.ToUpper(name)
}
