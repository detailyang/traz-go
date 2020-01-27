package traz

import (
	"testing"
	"time"

	"github.com/detailyang/size"
	"github.com/stretchr/testify/require"
)

func TestYAMLtraz(t *testing.T) {
	type MYYAML struct {
		A string
		B int
		C Duration
		D Size
	}
	x := []byte(`
a: "1"
b: 2
c: "5s"
d: "10KB"
`)
	tl := NewYAMLtraz(string(x))
	ml := &MYYAML{}
	err := tl.Apply(ml)
	require.Nil(t, err)
	require.Equal(t, "1", ml.A)
	require.Equal(t, 2, ml.B)
	require.Equal(t, Duration{5 * time.Second}, ml.C)
	require.Equal(t, Size{10 * size.KB}, ml.D)
}
