package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/db"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/rabbitmq"
)

var (
	cfgFile string
	//nolint:unused
	sigs = make(chan os.Signal, 1)
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

//nolint:deadcode,unused
func main() {
	flag.Parse()

	err := config.InitSchedulerConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	var st calendar.StorageInterface
	if config.Conf.SQL {
		st, err = db.NewSQLDatabase()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		st, err = db.NewInMemDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("connected to db")

	p, err := rabbitmq.NewPublisher()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to rabbit")

	app := scheduler.NewScheduler(st, p)
	defer func() {
		err := app.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("created scheduler")

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	errs := make(chan error, 1)

	go func() {
		for {
			err := app.Scan()
			if err != nil {
				errs <- err
			}
			log.Printf("scan done")
			time.Sleep(time.Minute)
		}
	}()

	select {
	case <-sigs:
		signal.Stop(sigs)
		return
	case err = <-errs:
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}
