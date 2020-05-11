package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("default test", func(t *testing.T) {
		envs, _ := ReadDir("testdata/env/")
		require.Equal(t, `"hello"`, envs["HELLO"])
		require.Equal(t, "bar", envs["BAR"])
		require.Equal(t, `   foo
with new line`, envs["FOO"])
		require.Equal(t, "", envs["UNSET"])
	})

	t.Run("non-existing directory", func(t *testing.T) {
		_, err := ReadDir("qw/")
		require.NotEqual(t, nil, err)
	})

	t.Run("not directory", func(t *testing.T) {
		_, err := ReadDir("test.sh")
		require.NotEqual(t, nil, err)
	})
}
