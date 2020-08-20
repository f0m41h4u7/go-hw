package config

import (
	"net"
	"os"

	"github.com/spf13/viper"
)

type LogLevel string

var Conf CalendarConfig

type CalendarConfig struct {
	HTTPServer ServerConfiguration
	GRPCServer ServerConfiguration
	Logger     LogConfiguration
	SQL        bool
	Database   DBConfiguration
}

func (c CalendarConfig) validate() error {
	if _, err := net.LookupHost(c.HTTPServer.Host); err != nil {
		return ErrWrongServerHost
	}

	if _, err := net.LookupHost(c.GRPCServer.Host); err != nil {
		return ErrWrongServerHost
	}

	if c.SQL {
		if _, err := net.LookupHost(c.Database.Host); err != nil {
			return ErrWrongDBHost
		}
	}

	return nil
}

func InitCalendarConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.SetConfigName("configs/calendar_config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return ErrCannotReadConfig
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		return ErrCannotParseConfig
	}

	return Conf.validate()
}
