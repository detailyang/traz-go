package traz

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/detailyang/size"
	"github.com/fatih/structs"
)

type FieldSet func(field *structs.Field, v string) (error, bool)

var (
	fieldSetHooks []FieldSet
)

func RegisterFieldSet(fs FieldSet) {
	fieldSetHooks = append(fieldSetHooks, fs)
}

func fieldSet(field *structs.Field, v string) error {
	for _, fs := range fieldSetHooks {
		err, ok := fs(field, v)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	switch f := field.Value().(type) {
	case Duration:
		var err error
		f.Duration, err = time.ParseDuration(v)
		if err != nil {
			return err
		}
		return field.Set(f)
	case Size:
		var err error
		f.Unit, err = size.Parse(v)
		if err != nil {
			return err
		}
		return field.Set(f)
	case flag.Value:
		if v := reflect.ValueOf(field.Value()); v.IsNil() {
			typ := v.Type()
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}

			if err := field.Set(reflect.New(typ).Interface()); err != nil {
				return err
			}

			f = field.Value().(flag.Value)
		}

		return f.Set(v)
	}

	// TODO: add support for other types
	switch field.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}

		if err := field.Set(val); err != nil {
			return err
		}
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if err := field.Set(i); err != nil {
			return err
		}
	case reflect.String:
		if err := field.Set(v); err != nil {
			return err
		}
	case reflect.Slice:
		switch t := field.Value().(type) {
		case []string:
			if err := field.Set(strings.Split(v, ",")); err != nil {
				return err
			}
		case []int:
			var list []int
			for _, in := range strings.Split(v, ",") {
				i, err := strconv.Atoi(in)
				if err != nil {
					return err
				}

				list = append(list, i)
			}

			if err := field.Set(list); err != nil {
				return err
			}
		default:
			return fmt.Errorf("traz: field '%s' of type slice is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		if err := field.Set(f); err != nil {
			return err
		}
	case reflect.Int64:
		switch t := field.Value().(type) {
		case size.Unit:
			d, err := size.Parse(v)
			if err != nil {
				return err
			}

			if err := field.Set(d); err != nil {
				return err
			}

		case time.Duration:
			d, err := time.ParseDuration(v)
			if err != nil {
				return err
			}

			if err := field.Set(d); err != nil {
				return err
			}
		case int64:
			p, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return err
			}

			if err := field.Set(p); err != nil {
				return err
			}
		default:
			return fmt.Errorf("traz: field '%s' of type int64 is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}

	default:
		return fmt.Errorf("traz: field '%s' has unsupported type: %s", field.Name(), field.Kind())
	}

	return nil
}
