package traz

import (
	"testing"
	"time"

	"github.com/detailyang/size"
	"github.com/stretchr/testify/require"
)

func TestTOMLtraz(t *testing.T) {
	type MYToml struct {
		A string
		B int
		C Duration
		D Size
	}
	x := []byte(`
a = "1"
b = 2
c = "5s"
d = "10KB"
`)
	tl := NewTOMLtraz(string(x))
	ml := &MYToml{}
	err := tl.Apply(ml)
	require.Nil(t, err)
	require.Equal(t, "1", ml.A)
	require.Equal(t, 2, ml.B)
	require.Equal(t, Duration{5 * time.Second}, ml.C)
	require.Equal(t, Size{10 * size.KB}, ml.D)
}
