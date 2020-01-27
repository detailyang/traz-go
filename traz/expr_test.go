package traz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExprtraz(t *testing.T) {
	type Foo struct {
		A string `expr:"len(A) > 0 && len(A) < 16"`
	}

	vz := &Exprtraz{}
	foo := &Foo{}
	err := vz.Apply(foo)
	require.NotNil(t, err)
	foo.A = "abcd"
	err = vz.Apply(foo)
	require.Nil(t, err)

	type Bar struct {
		B int `expr:"B > 0 && B < 16"`
	}
	vz = &Exprtraz{}
	bar := &Bar{}
	err = vz.Apply(bar)
	require.NotNil(t, err)
	bar.B = 2
	err = vz.Apply(foo)
	require.Nil(t, err)
}
