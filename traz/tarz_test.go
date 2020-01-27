package traz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiTraz(t *testing.T) {
	type Foo struct{}
	foo := &Foo{}
	mz := NewMultiTraz()
	mz.Append(TransformerFunc(func(s interface{}) error {
		require.Equal(t, foo, s)
		return nil
	}))
	mz.Append(TransformerFunc(func(s interface{}) error {
		require.Equal(t, foo, s)
		return nil
	}))
	mz.Append(TransformerFunc(func(s interface{}) error {
		require.Equal(t, foo, s)
		return nil
	}))
	err := mz.Apply(foo)
	require.Nil(t, err)
}
