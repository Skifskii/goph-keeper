package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Skifskii/goph-keeper/internal/config"
	"github.com/Skifskii/goph-keeper/internal/repository/postgres"
	"github.com/Skifskii/goph-keeper/internal/service"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv"
	"github.com/Skifskii/goph-keeper/pkg/crypto"
)

func Run() error {
	// config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

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

	// adapters
	cryp, err := crypto.New([]byte(cfg.MasterKey))
	if err != nil {
		return fmt.Errorf("failed to initialize crypto: %w", err)
	}

	// services
	serv, err := service.New(
		repo,
		repo,
		cryp,
		[]byte(cfg.MasterKey),
		cfg.SecretKey,
		cfg.JWTTokenTTL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// transport
	httpServer := httpserv.New(log, cfg.HTTP.Address, serv.Auth, serv.Secret, serv.Secret)

	return httpServer.Run()
}
