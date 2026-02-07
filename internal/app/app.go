package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	pgRepo, err := postgres.New(log, cfg.DatabaseDSN)
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
		pgRepo,
		pgRepo,
		cryp,
		[]byte(cfg.MasterKey),
		cfg.SecretKey,
		cfg.JWTTokenTTL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// transport
	httpServer := httpserv.New(
		log,
		cfg.HTTP.Address,
		serv.Auth,
		serv.Secret,
		serv.Secret,
		serv.Secret,
		serv.Secret,
		serv.Secret,
	)
	go httpServer.MustRun()

	// graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	sign := <-signalChan
	log.Info("signal received, starting to shut down", slog.Any("signal", sign))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		log.Error("failed to stop http server", slog.Any("error", err))
	}
	if err := pgRepo.Stop(); err != nil {
		log.Error("failed to stop postgres repo", slog.Any("error", err))
	}

	return nil
}
