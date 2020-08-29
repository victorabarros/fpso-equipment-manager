package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config summarises environment variables.
type Config struct {
	LogLevel string `mapstructure:"log_level"`
}

// Load return all environment variables loaded.
func Load() (cfg Config, err error) {
	logrus.Debug("Loading enviromnts variables.")

	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.AutomaticEnv()
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}
	return
}
