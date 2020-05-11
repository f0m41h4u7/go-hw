package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint
	cmdExec.Env = os.Environ()
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	for envVar, val := range env {
		newVar := envVar + "=" + val
		cmdExec.Env = append(cmdExec.Env, newVar)
	}

	err := cmdExec.Start()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmdExec.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			log.Fatal(err)
		}
	}

	return 0
}
