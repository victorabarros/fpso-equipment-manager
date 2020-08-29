package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/victorabarros/challenge-modec/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Error in load Enviromnts variables.")
	}

	fmt.Println(cfg)
}
