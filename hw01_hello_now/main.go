package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"time"
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
