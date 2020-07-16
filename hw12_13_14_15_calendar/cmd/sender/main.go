package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/rabbitmq"
)

var (
	cfgFile string
	sigs    = make(chan os.Signal, 1)
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

func main() {
	flag.Parse()

	err := config.InitSenderConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	c := rabbitmq.NewConsumer()

	app := sender.NewSender(c)
	defer func() {
		err := app.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	errs := make(chan error, 1)

	go func() {
		errs <- app.Listen()
	}()

	select {
	case <-sigs:
		signal.Stop(sigs)
		return
	case err = <-errs:
		if err != nil {
			log.Fatal(err)
		}
	}
}
