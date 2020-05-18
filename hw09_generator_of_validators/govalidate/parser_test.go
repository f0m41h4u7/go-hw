package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsingTags(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		err := parseTag("`validate:\"len:42\"`", "float32")
		require.NotNil(t, err)

		err = parseTag("`validate:\"len:42\"`", "string")
		require.Nil(t, err)
		require.Equal(t, "42", fLen)
	})

	t.Run("empty tag", func(t *testing.T) {
		err := parseTag("`validate:`", "string")
		require.Nil(t, err)
	})

	t.Run("variants", func(t *testing.T) {
		err := parseTag("`validate:\"in:14333,12\"`", "int")
		require.Nil(t, err)
		require.Equal(t, "14333,12", fIn)

		err = parseTag("`validate:\"in:14333,12\"`", "string")
		require.Nil(t, err)
		require.Equal(t, "\"14333\",\"12\"", fIn)

		err = parseTag("`validate:\"in:twelve,12\"`", "string")
		require.Nil(t, err)
		require.Equal(t, "\"twelve\",\"12\"", fIn)
	})
}
