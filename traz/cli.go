package traz

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/spf13/pflag"
)

type Clitraz struct {
	name   string
	prefix string
	args   []string
	fs     *pflag.FlagSet
}

func NewClitraz(name, prefix string, args []string) *Clitraz {
	fs := pflag.NewFlagSet(name, pflag.ExitOnError)
	ct := &Clitraz{
		name:   name,
		prefix: prefix,
		args:   args,
		fs:     fs,
	}
	return ct
}

func (ct *Clitraz) GenUsage(w io.Writer, s interface{}) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Cli Usage of %s:\n", ct.name)
		ct.fs.PrintDefaults()
		fmt.Fprintf(w, "\nEnv Usage of %s:\n", ct.name)
		e := NewEnvtraz(ct.prefix)
		// nolint
		e.String(os.Stderr, s)
		fmt.Println("")
	}
}

func (ct *Clitraz) Apply(s interface{}) error {
	ct.fs.Usage = ct.GenUsage(os.Stderr, s)

	prefix := ct.genprefix(ct.prefix, structs.Name(s))
	for _, field := range structs.Fields(s) {
		if err := ct.apply(prefix, field); err != nil {
			return err
		}
	}

	return ct.fs.Parse(ct.args)
}

func (ct *Clitraz) apply(prefix string, field *structs.Field) error {
	prefix = ct.genprefix(prefix, field.Name())
	switch field.Kind() {
	case reflect.Struct:
		switch field.Value().(type) {
		case Duration, Size:
			ct.fs.Var(newFieldValue(field), flagName(prefix), ct.flagUsage(prefix, field))
			return nil
		}

		for _, ff := range field.Fields() {
			if err := ct.apply(prefix, ff); err != nil {
				return err
			}
		}

	default:
		if field.IsExported() {
			ct.fs.Var(newFieldValue(field), flagName(prefix), ct.flagUsage(prefix, field))
		}
	}
	return nil
}

func (ct *Clitraz) flagUsage(fieldName string, field *structs.Field) string {
	usage := field.Tag("cli")
	if usage != "" {
		return usage
	}

	return fmt.Sprintf("Change value of %q.", fieldName)
}

func (ct *Clitraz) genprefix(prefix string, name string) string {
	if prefix == "" {
		return strings.ToUpper(name)
	}
	return strings.ToUpper(prefix) + "." + strings.ToUpper(name)
}

// fieldValue satisfies the flag.Value and flag.Getter interfaces
type fieldValue struct {
	field *structs.Field
}

func newFieldValue(f *structs.Field) *fieldValue {
	return &fieldValue{
		field: f,
	}
}

func (f *fieldValue) Set(val string) error {
	return fieldSet(f.field, val)
}

func (f *fieldValue) String() string {
	if f.IsZero() {
		return ""
	}

	return fmt.Sprintf("%v", f.field.Value())
}

func (f *fieldValue) Type() string {
	kind := f.field.Kind().String()
	switch f.field.Value().(type) {
	case Duration:
		kind = "duration"
	case Size:
		kind = "size"
	}
	return kind
}

func (f *fieldValue) Get() interface{} {
	if f.IsZero() {
		return nil
	}

	return f.field.Value()
}

func (f *fieldValue) IsZero() bool {
	return f.field == nil
}

// This is an unexported interface, be careful about it.
// https://code.google.com/p/go/source/browse/src/pkg/flag/flag.go?name=release#101
func (f *fieldValue) IsBoolFlag() bool {
	return f.field.Kind() == reflect.Bool
}

func flagName(name string) string { return strings.ToLower(name) }
