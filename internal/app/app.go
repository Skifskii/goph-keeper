package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Skifskii/goph-keeper/internal/config"
	"github.com/Skifskii/goph-keeper/internal/repository/postgres"
	"github.com/Skifskii/goph-keeper/internal/service"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv"
)

func Run() error {
	// config
	cfg := config.New()

	// logger
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	// repository
	repo, err := postgres.New(log, cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize repo: %w", err)
	}

	// services
	serv, err := service.New(
		repo,
		[]byte(cfg.MasterKey),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// transport
	httpServer := httpserv.New(log, cfg.HTTP.Address, serv.Secret)

	return httpServer.Run()
}
