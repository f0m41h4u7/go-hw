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
	timeout          time.Duration
	ErrNotEnoughArgs = errors.New("not enough arguments, should be 3 at least")
	sigs             = make(chan os.Signal, 1)
)

const (
	minLenArgs = 3
	maxLenArgs = 4
)

func init() {
	flag.DurationVar(&timeout, "timeout", 0, "connection timeout")
}

func main() {
	flag.Parse()
	if (len(os.Args) < minLenArgs) || (len(os.Args) > maxLenArgs) {
		log.Fatal(ErrNotEnoughArgs)
	}
	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]

	c := NewTelnetClient(
		net.JoinHostPort(host, port),
		timeout,
		os.Stdin,
		os.Stdout,
	)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err := c.Connect()
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

	select {
	case <-sigs:
		signal.Stop(sigs)
		return
	case err = <-errs:
		if err != nil {
			log.Fatal(err)
		}
		errLog.Println("...EOF")
		return
	}
}
