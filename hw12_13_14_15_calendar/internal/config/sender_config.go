package config

import (
	"net"
	"os"

	"github.com/spf13/viper"
)

var SendConf SenderConfig

type SenderConfig struct {
	Rabbit ServerConfiguration
}

func (c SenderConfig) validate() error {
	if _, err := net.LookupHost(c.Rabbit.Host); err != nil {
		return ErrWrongRMQHost
	}

	return nil
}

func InitSenderConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.SetConfigName("configs/sender_config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return ErrCannotReadConfig
	}

	if err := viper.Unmarshal(&SendConf); err != nil {
		return ErrCannotParseConfig
	}

	return SendConf.validate()
}
