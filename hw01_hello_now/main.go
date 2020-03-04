package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal(err)
	}
	localTime := time.Now()

	fmt.Printf("current time: %v\n", localTime)
	fmt.Printf("exact time: %v\n", exactTime)
}
