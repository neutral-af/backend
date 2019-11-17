package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Environment     string
	CloverlyAPIKey  string `mapstructure:"CLOVERLY_API_KEY"`
	HoneycombAPIKey string `mapstructure:"HONEYCOMB_API_KEY"`
	StripeSecretKey string `mapstructure:"STRIPE_SECRET_KEY"`
}

// C is the object containing config values
var C config

const branchEnvVar = "NOW_GITHUB_COMMIT_REF"

func init() {
	C = New()
}

func New() config {
	var c config

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("dotenv")
	}

	ref := os.Getenv(branchEnvVar)
	switch {
	case ref == "master":
		c.Environment = "prod"
	case ref != "":
		c.Environment = "staging"
	default:
		c.Environment = "dev"
	}

	viper.SetEnvPrefix(c.Environment)

	viper.BindEnv("CLOVERLY_API_KEY")
	viper.BindEnv("HONEYCOMB_API_KEY")
	viper.BindEnv("STRIPE_SECRET_KEY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	viper.Unmarshal(&c)
	return c
}
