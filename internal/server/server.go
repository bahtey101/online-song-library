package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"online-song-library/internal/bootstrap"
	"online-song-library/internal/clients/infoservice"
	"online-song-library/internal/config"
	"online-song-library/internal/handler"
	"online-song-library/internal/repository/songrepository"
	"online-song-library/internal/service"

	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())

	logrus.SetLevel(logrus.Level(cfg.LogLevel))

	pgConnPool, err := bootstrap.InitDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect postgres %s, %v", cfg.PgDSN, err)
	}

	songRepository := songrepository.NewSongRepository(pgConnPool)
	client := infoservice.NewMusicInfoClient(cfg)
	service := service.NewService(songRepository, client)
	handler := handler.NewHandler(service)

	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler.InitRoutes(),
		MaxHeaderBytes: 1 << 28, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		logrus.Infof("Starting HTTP server on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error service http server %v", err)
		}
	}()

	gracefulShotdown(ctx, server, cancel)

	return nil
}

func gracefulShotdown(ctx context.Context, s *http.Server, cancel context.CancelFunc) {
	const waitTime = 5 * time.Second // waiting time before closing all connections

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(ch)

	sig := <-ch
	logrus.Infof("Received shutdown signal: %v. Initiating graceful shutdown...", sig)

	if err := s.Shutdown(ctx); err != nil {
		logrus.Errorf("Error shutting down server: %v", err)
	}

	cancel()
	time.Sleep(waitTime)
	logrus.Info("Graceful shutdown completed.")
}
