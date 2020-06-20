package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitConfig("../../../tests/testdata/config.json")
		require.Nil(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		err := InitConfig("../../../tests/testdata/wrong_host.json")
		require.Equal(t, ErrWrongServerHost, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := InitConfig("config.json")
		require.Equal(t, ErrWrongConfig, err)
	})

	t.Run("wrong config structure", func(t *testing.T) {
		err := InitConfig("../../../tests/testdata/bad_structure.json")
		require.Equal(t, ErrWrongConfig, err)
	})
}
