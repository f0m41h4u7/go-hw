package tests

import (
	"os"
	"testing"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/config"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	pwd, _ := os.Getwd()

	t.Run("simple", func(t *testing.T) {
		err := config.InitConfig(pwd + "/test_configs/config.json")
		require.Nil(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		err := config.InitConfig(pwd + "/test_configs/wrong_host.json")
		require.Equal(t, config.ErrWrongServerHost, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := config.InitConfig(pwd + "config.json")
		require.Equal(t, config.ErrWrongConfig, err)
	})

	t.Run("wrong config structure", func(t *testing.T) {
		err := config.InitConfig(pwd + "/test_configs/bad_structure.json")
		require.Equal(t, config.ErrWrongConfig, err)
	})
}
