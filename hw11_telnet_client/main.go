package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout          string
	ErrNotEnoughArgs = errors.New("not enough arguments, should be 3 at least")
	sigs             = make(chan os.Signal, 1)
)

const minLenArgs = 3

func init() {
	flag.StringVar(&timeout, "timeout", "0", "connection timeout")
}

func main() {
	flag.Parse()
	if len(os.Args) < minLenArgs {
		log.Fatal(ErrNotEnoughArgs)
	}
	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]

	validTime, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal(err)
	}

	c := NewTelnetClient(
		net.JoinHostPort(host, port),
		validTime,
		os.Stdin,
		os.Stdout,
	)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err = c.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	errs := make(chan error, 1)
	go func() { errs <- c.Send() }()
	go func() { errs <- c.Receive() }()

	for {
		select {
		case <-sigs:
			return
		case err = <-errs:
			if err != nil {
				log.Fatal(err)
			}
			errLog.Println("...EOF")
			return
		}
	}
}
