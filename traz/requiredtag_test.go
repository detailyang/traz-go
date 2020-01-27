package traz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Bar struct {
	A string `required:"true"`
}

func TestRequiredTag(t *testing.T) {
	c := NewRequiredTagtraz(DefaultRequiredTagName)
	bar := &Bar{}
	err := c.Apply(bar)
	require.NotNil(t, err)
	bar.A = "a"
	err = c.Apply(bar)
	require.Nil(t, err)
}
