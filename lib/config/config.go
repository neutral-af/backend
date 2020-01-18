package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Environment string

func (e Environment) ToString() string {
	return string(e)
}

const (
	EnvironmentProd    = "prod"
	EnvironmentStaging = "staging"
	EnvironmentDev     = "dev"
)

type config struct {
	Environment               Environment
	CloverlyAPIKey            string `mapstructure:"CLOVERLY_API_KEY"`
	DigitalHumaniEnterpriseID string `mapstructure:"DIGITALHUMANI_ENTERPRISE_ID"`
	FlightStatsAppID          string `mapstructure:"FLIGHTSTATS_APP_ID"`
	FlightStatsAppKey         string `mapstructure:"FLIGHTSTATS_APP_KEY"`
	HoneycombAPIKey           string `mapstructure:"HONEYCOMB_API_KEY"`
	StripeSecretKey           string `mapstructure:"STRIPE_SECRET_KEY"`
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
		c.Environment = EnvironmentProd
	case ref != "":
		c.Environment = EnvironmentStaging
	default:
		c.Environment = EnvironmentDev
	}

	viper.SetEnvPrefix(c.Environment.ToString())

	viper.BindEnv("CLOVERLY_API_KEY")
	viper.BindEnv("DIGITALHUMANI_ENTERPRISE_ID")
	viper.BindEnv("HONEYCOMB_API_KEY")
	viper.BindEnv("STRIPE_SECRET_KEY")
	viper.BindEnv("FLIGHTSTATS_APP_ID")
	viper.BindEnv("FLIGHTSTATS_APP_KEY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	viper.Unmarshal(&c)
	return c
}
