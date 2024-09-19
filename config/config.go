package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	Port    int    `mapstructure:"port"`
}

var cfg Config

func init() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
}

func GetConfig() Config {
	return cfg
}
