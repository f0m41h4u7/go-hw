package main

import (
	"flag"
	"log"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/http"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/logger"
)

var (
	cfgFile string
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

	http.StartServer()
}
