package main

import (
	"fmt"
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
		log.Fatal(fmt.Errorf(""))
	}
}
