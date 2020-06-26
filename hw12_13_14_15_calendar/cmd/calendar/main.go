package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/api/grpcapi"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/api/httpapi"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/db"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
)

var (
	cfgFile string
	app     calendar.Calendar
	sigs    = make(chan os.Signal, 1)
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

func main() {
	flag.Parse()

	err := config.InitConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	err = logger.InitLogger()
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

	app = calendar.Calendar{
		Storage: st,
	}

	http := httpapi.InitServer(&app)
	defer func() {
		if err := http.Stop(); err != nil {
			zap.L().Error("http server stop with error", zap.Error(err))
			return
		}
	}()
	grpc := grpcapi.InitServer(&app)
	defer grpc.Stop()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	errs := make(chan error, 1)
	go func() { errs <- http.Start() }()
	go func() { errs <- grpc.Start() }()

	select {
	case <-sigs:
		signal.Stop(sigs)
		return
	case err = <-errs:
		if err != nil {
			zap.L().Error("server exited with error", zap.Error(err))
		}
		return
	}
}
