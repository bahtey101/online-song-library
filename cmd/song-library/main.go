package main

import (
	"online-song-library/internal/config"
	"online-song-library/internal/server"

	"github.com/caarlos0/env/v11"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := new(config.Config)

	if err := env.Parse(cfg); err != nil {
		logrus.Fatal("failed to retrieve env variables: ", err)
	}

	if err := server.Run(cfg); err != nil {
		logrus.Fatal("error running service: ", err)
	}
}
