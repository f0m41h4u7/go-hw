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

var (
	cfgFile string
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
		conn, err := db.InitSQLConnection()
		if err != nil {
			log.Fatal(err)
		}
		st = db.NewSQLDatabase(conn)
	} else {
		st, err = db.NewInMemDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}

	calendar.TheCalendar = calendar.Calendar{
		Storage: st,
	}

	httpserver.StartServer()
}
