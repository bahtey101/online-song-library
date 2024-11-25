package main

import (
	"online-song-library/internal/config"
	"online-song-library/internal/server"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "online-song-library/docs"
)

// @title Song library
// @version 0.0.1
// @description Online song library on Go
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://github.com/bahtey101/online-song-library/blob/main/LICENSE

// @host localhost:8080
// @BasePath /songs
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
