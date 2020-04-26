package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy /dev/urandom", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})
	t.Run("file doesn't exist", func(t *testing.T) {
		err := Copy("test.txt", "/tmp", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("copy directory", func(t *testing.T) {
		err := Copy("/tmp", "test.txt", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})
	t.Run("copy to directory", func(t *testing.T) {
		_, _ = os.Create("test.txt")
		defer os.Remove("test.txt")
		err := Copy("test.txt", "/tmp", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})
	t.Run("permissions", func(t *testing.T) {
		_, _ = os.Create("test1.txt")
		defer os.Remove("test1.txt")

		_, _ = os.Create("test2.txt")
		defer os.Remove("test2.txt")

		_ = os.Chmod("test2.txt", 0444)

		err := Copy("test1.txt", "test2.txt", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("copy file to itself", func(t *testing.T) {
		f, _ := os.Create("test.txt")
		defer os.Remove("test.txt")
		_, _ = f.WriteString("test")

		err := Copy("test.txt", "test.txt", 0, 0)
		require.NotNil(t, err)
	})
}
