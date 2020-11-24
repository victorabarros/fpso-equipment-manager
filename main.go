package main

import (
	"github.com/sirupsen/logrus"
	"github.com/victorabarros/fpso-equipment-manager/app/server"
	"github.com/victorabarros/fpso-equipment-manager/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Error in load Enviromnts variables.")
	}

	loglvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatalf(
			"Error in set log level %s.", cfg.LogLevel)
	}
	logrus.SetLevel(loglvl)

	server.Run("8092")
}
