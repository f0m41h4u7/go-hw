package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/db"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/rabbitmq"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

func main() {
	flag.Parse()

	err := config.InitSchedulerConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	var st scheduler.Storage
	if config.SchedConf.SQL {
		st, err = db.NewSQLDatabase(config.SchedConf.Database)
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

	app := scheduler.NewScheduler(st, rabbitmq.NewPublisher())
	defer func() {
		err := app.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("created scheduler")

	sigs := make(chan os.Signal, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(time.Duration(config.SchedConf.Interval) * time.Millisecond)
	defer cancel()

	go func() {
		for {
			select {
			case <-sigs:
				signal.Stop(sigs)
				return
			case <-ticker.C:
				app.Scan()
			case <-ctx.Done():
				return
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigs
}
