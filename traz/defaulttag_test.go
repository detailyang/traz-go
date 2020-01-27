package traz

import (
	"testing"
	"time"

	"github.com/detailyang/size"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	A string        `default:"abcd"`
	B string        `default:"defg"`
	C time.Duration `default:"5s"`
	D int           `default:"100"`
	E size.Unit     `default:"2KB"`
	F Size          `default:"8KB"`
	G Duration      `default:"6s"`
}

func TestDefaultTag(t *testing.T) {
	c := NewDefaultTagtraz(DefaultDefaultTagName)
	foo := &Foo{}
	err := c.Apply(foo)
	require.Nil(t, err)
	require.Equal(t, "abcd", foo.A)
	require.Equal(t, "defg", foo.B)
	require.Equal(t, 5*time.Second, foo.C)
	require.Equal(t, 100, foo.D)
	require.Equal(t, 2*size.KB, foo.E)
	require.Equal(t, Size{8 * size.KB}, foo.F)
	require.Equal(t, Duration{6 * time.Second}, foo.G)
}
