package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	Conf           Config
	ErrWrongConfig = errors.New("cannot parse config file")
)

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
		viper.SetConfigName("config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(ErrWrongConfig)
	}

	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatal(ErrWrongConfig)
	}
}
