package traz

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/fatih/structs"
)

// Exprtraz executes the expr from struct tag via https://github.com/antonmedv/expr
type Exprtraz struct {
	env map[string]interface{}
}

func NewExprtraz() *Exprtraz {
	return &Exprtraz{
		env: make(map[string]interface{}, 32),
	}
}

func (vt *Exprtraz) RegisterEnv(name string, v interface{}) {
	vt.env[name] = v
}

func (vt *Exprtraz) RegisterEnvMap(m map[string]interface{}) {
	for k, v := range m {
		vt.env[k] = v
	}
}

func (vt *Exprtraz) Apply(s interface{}) error {
	prefix := vt.genprefix("", structs.Name(s))
	strct := structs.New(s)
	for _, field := range structs.Fields(s) {
		if err := vt.apply(prefix, strct, field); err != nil {
			return err
		}
	}
	return nil
}
func (vt *Exprtraz) apply(prefix string, strct *structs.Struct,
	field *structs.Field) error {
	prefix = vt.genprefix(prefix, field.Name())
	switch field.Kind() {
	case reflect.Struct:
		strct = structs.New(field.Value())
		for _, f := range field.Fields() {
			if err := vt.apply(prefix, strct, f); err != nil {
				return err
			}
		}
	default:
		return vt.Validate(prefix, strct, field)
	}

	return nil
}

func (vt *Exprtraz) Validate(prefix string, strct *structs.Struct, field *structs.Field) error {
	src := field.Tag("expr")
	if src == "" {
		return nil
	}

	environment := make(map[string]interface{}, 16)
	strct.FillMap(environment)

	for k, v := range vt.env {
		environment[k] = v
	}

	program, err := expr.Compile(src, expr.Env(environment))
	if err != nil {
		return err
	}

	output, err := expr.Run(program, environment)
	if err != nil {
		return err
	}

	switch t := output.(type) {
	case bool:
		if !t {
			return fmt.Errorf("%q return false", src)
		}
	default:
		return fmt.Errorf("%q return non-bool type", src)
	}

	return nil
}

func (vt *Exprtraz) genprefix(prefix string, name string) string {
	if prefix == "" {
		return strings.ToUpper(name)
	}
	return strings.ToUpper(prefix) + "." + strings.ToUpper(name)
}

var (
	BuiltInMap = map[string]interface{}{
		"IsIP4":    IsIP4,
		"IsIP6":    IsIP6,
		"IsIP":     IsIP,
		"IsPort":   IsPort,
		"IsIPPort": IsIPPort,
	}
)

func IsIP4(s string) bool {
	ip := net.ParseIP(s)
	return ip.To4() != nil
}

func IsIP6(s string) bool {
	ip := net.ParseIP(s)
	return ip.To16() != nil
}

func IsIP(s string) bool {
	ip := net.ParseIP(s)
	if ip.To4() != nil || ip.To16() != nil {
		return true
	}
	return false
}

func IsPort(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if n >= 65535 || n < 0 {
		return false
	}
	return true
}

func IsIPPort(s string) bool {
	h, p, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}

	return IsPort(p) && IsIP(h)
}
