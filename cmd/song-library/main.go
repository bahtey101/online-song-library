package main

import (
	"online-song-library/internal/config"
	"online-song-library/internal/server"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := new(config.Config)

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("No .env file found")
	}

	if err := env.Parse(cfg); err != nil {
		logrus.Fatal("Failed to retrieve env variables: ", err)
	}

	if err := server.Run(cfg); err != nil {
		logrus.Fatal("Error running service: ", err)
	}
}
