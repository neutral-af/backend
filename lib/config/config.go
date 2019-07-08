package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	CloverlyAPIKey string `mapstructure:"CLOVERLY_API_KEY"`
}

// C is the object containing config values
var C config

func init() {
	C = New()
}

func New() config {
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("dotenv")
	}

	viper.BindEnv("CLOVERLY_API_KEY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	var c config
	viper.Unmarshal(&c)
	return c
}
