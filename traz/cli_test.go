package traz

import (
	"testing"
	"time"

	"github.com/detailyang/size"
	"github.com/stretchr/testify/require"
)

func TestClitraz(t *testing.T) {
	type Bar struct {
		C string
		D Duration
		E Size
	}

	type Foo struct {
		A   string
		B   int
		Bar Bar
	}

	ct := NewClitraz("test", "", []string{"--foo.a=abcd", "--foo.b=1",
		"--foo.bar.c=ee", "--foo.bar.d=5s", "--foo.bar.e=5KB"})
	foo := &Foo{}
	err := ct.Apply(foo)
	require.Nil(t, err)
	require.Equal(t, "abcd", foo.A)
	require.Equal(t, 1, foo.B)
	require.Equal(t, "ee", foo.Bar.C)
	require.Equal(t, Duration{5 * time.Second}, foo.Bar.D)
	require.Equal(t, Size{5 * size.KB}, foo.Bar.E)
}
