package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	QWCURL        string `mapstructure:"QWCURL" validate:"required"`
	QWCTOKEN      string `mapstructure:"QWCTOKEN" validate:"required"`
	MEDIAS3LAMBDA string `mapstructure:"MEDIAS3LAMBDA" validate:"required"`
	XPUSERNAME    string `mapstructure:"XPUSERNAME" validate:"required"`
	XPPASSWORD    string `mapstructure:"XPPASSWORD" validate:"required"`
	DOLBYUSERNAME string `mapstructure:"DOLBYUSERNAME" validate:"required"`
	DOLBYPASSWORD string `mapstructure:"DOLBYPASSWORD" validate:"required"`
}

func LoadConfig() *EnvConfig {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	var envConfig EnvConfig

	if err := viper.Unmarshal(&envConfig); err != nil {
		panic(err)
	}

	return &envConfig
}
