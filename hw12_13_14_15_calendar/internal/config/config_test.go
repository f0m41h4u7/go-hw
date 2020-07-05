package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalendarConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitCalendarConfig("../../tests/testdata/calendar_config.json")
		require.Nil(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		err := InitCalendarConfig("../../tests/testdata/calendar_wrong_host.json")
		require.Equal(t, ErrWrongServerHost, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := InitCalendarConfig("config.json")
		require.Equal(t, ErrCannotReadConfig, err)
	})

	t.Run("wrong config structure", func(t *testing.T) {
		err := InitCalendarConfig("../../tests/testdata/calendar_bad_structure.json")
		require.Equal(t, ErrCannotParseConfig, err)
	})
}

func TestSchedulerConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitSchedulerConfig("../../tests/testdata/scheduler_config.json")
		require.Nil(t, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := InitSchedulerConfig("config.json")
		require.Equal(t, ErrCannotReadConfig, err)
	})
}

func TestSenderConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitSenderConfig("../../tests/testdata/sender_config.json")
		require.Nil(t, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := InitSenderConfig("config.json")
		require.Equal(t, ErrCannotReadConfig, err)
	})
}
