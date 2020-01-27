package traz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapPasstraz(t *testing.T) {
	type Bar struct {
		A string
	}
	type Foo struct {
		M map[string]interface{}
	}

	m := NewMapPasstraz()
	m.Append(func(prefix string, name string, m interface{}) (interface{}, bool, error) {
		if prefix != "FOO.M.A" {
			return nil, false, nil
		}

		return &Bar{
			A: "abcd",
		}, true, nil
	})

	foo := &Foo{
		M: map[string]interface{}{
			"A": "1234",
		},
	}
	err := m.Apply(foo)
	require.Nil(t, err)
	require.Equal(t, &Bar{A: "abcd"}, foo.M["A"])
}
