package traz

import (
	"bytes"
	"testing"
	"time"

	"github.com/detailyang/size"

	"github.com/stretchr/testify/require"
)

func TestEnvtraz(t *testing.T) {
	type Bar struct {
		C string
		D Duration
		E Size
	}

	type Foo struct {
		A   string
		B   int
		Bar Bar
		M   map[string]interface{}
	}

	et := NewEnvtraz("")
	et.getter = EnvGetterFunc(func(key string) string {
		m := map[string]string{
			"FOO.A":          "1",
			"FOO.B":          "2",
			"FOO.BAR.C":      "3",
			"FOO.BAR.D":      "10s",
			"FOO.BAR.E":      "10KB",
			"FOO.M.STRUCT.A": "abcd",
			"FOO.M.STRUCT.B": "5s",
		}
		return m[key]
	})

	foo := &Foo{
		M: map[string]interface{}{
			"struct": &struct {
				A string
				B Duration
			}{},
		},
	}
	err := et.Apply(foo)
	require.Nil(t, err)
	require.Equal(t, "1", foo.A)
	require.Equal(t, 2, foo.B)
	require.Equal(t, "3", foo.Bar.C)
	require.Equal(t, Duration{10 * time.Second}, foo.Bar.D)
	require.Equal(t, Size{10 * size.KB}, foo.Bar.E)
	require.Equal(t, &struct {
		A string
		B Duration
	}{"abcd", Duration{5 * time.Second}}, foo.M["struct"])

	buf := bytes.NewBuffer(nil)
	err = et.String(buf, foo)
	require.Nil(t, err)
}
