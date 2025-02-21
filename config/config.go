package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration of the application.
// The values are read by viper from config file or env variables.
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddr          string        `mapstructure:"SERVER_ADDR"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// Load reads configuration from a file or env variables.
func Load(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
