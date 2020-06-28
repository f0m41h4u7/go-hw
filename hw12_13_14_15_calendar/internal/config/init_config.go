package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Conf                 Config
	ErrCannotReadConfig  = errors.New("cannot read config file")
	ErrCannotParseConfig = errors.New("cannot parse config file")
)

func InitConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
		viper.SetConfigName("configs/config")
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
