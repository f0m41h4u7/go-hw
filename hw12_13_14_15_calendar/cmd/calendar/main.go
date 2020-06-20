package cmd

import (
	"flag"
	"log"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/db"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar/httpserver"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/logger"
)

//nolint:unused
var (
	cfgFile string
	app     calendar.Calendar
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

//nolint:deadcode,unused
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

	httpserver.StartServer()
}
