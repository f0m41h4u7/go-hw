package config

import (
	"net"
	"os"

	"github.com/spf13/viper"
)

var SchedConf SchedulerConfig

type SchedulerConfig struct {
	Rabbit   ServerConfiguration
	Interval int64
	SQL      bool
	Database DBConfiguration
}

func (c SchedulerConfig) validate() error {
	if _, err := net.LookupHost(c.Rabbit.Host); err != nil {
		return ErrWrongRMQHost
	}

	if c.SQL {
		if _, err := net.LookupHost(c.Database.Host); err != nil {
			return ErrWrongDBHost
		}
	}
	return nil
}

func InitSchedulerConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.SetConfigName("configs/scheduler_config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return ErrCannotReadConfig
	}

	if err := viper.Unmarshal(&SchedConf); err != nil {
		return ErrCannotParseConfig
	}
	return SchedConf.validate()
}
