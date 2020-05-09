package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/echo2.sh"}
		envs := make(Environment)
		envs["TEST"] = "test"

		r, w, _ := os.Pipe()
		os.Stdout = w

		_ = RunCmd(cmd, envs)

		w.Close()
		out, _ := ioutil.ReadAll(r)

		require.Equal(t, "test", strings.TrimRight(string(out), " \t\n"))
	})

	t.Run("non-existing command", func(t *testing.T) {
		cmd := []string{"/bin/bash", "echo.sh"}
		envs := make(Environment)
		envs["TEST"] = "test"

		ret := RunCmd(cmd, envs)

		require.Equal(t, 127, ret)
	})

	t.Run("empty env", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/echo2.sh"}
		os.Setenv("TEST", "default")
		envs := make(Environment)

		r, w, _ := os.Pipe()
		os.Stdout = w

		_ = RunCmd(cmd, envs)

		w.Close()
		out, _ := ioutil.ReadAll(r)

		require.Equal(t, "default", strings.TrimRight(string(out), " \t\n"))
	})
}
