package main

import (
	"log"
	"os"
)

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(os.Args[2:], envs)
	if returnCode != 0 {
		log.Printf("command exited with code %d\n", returnCode)
		log.Fatal()
	}
}
