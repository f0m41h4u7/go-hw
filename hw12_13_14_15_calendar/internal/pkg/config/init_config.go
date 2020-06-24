package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Conf           Config
	ErrWrongConfig = errors.New("cannot parse config file")
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
		return ErrWrongConfig
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		return ErrWrongConfig
	}
	return Conf.validate()
}
