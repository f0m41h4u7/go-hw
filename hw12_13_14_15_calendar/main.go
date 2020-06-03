package main

import (
	"flag"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/http"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/logger"
)

var (
	cfgFile string
)

const minLenArgs = 2

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

func main() {
	flag.Parse()

	config.InitConfig(cfgFile)
	logger.InitLogger()
	http.StartServer()
}
