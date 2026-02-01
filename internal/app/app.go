package app

import (
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
	_ = postgres.New()

	// services
	_ = service.New()

	// transport
	httpServer := httpserv.New(log, cfg.HTTP.Address)

	return httpServer.Run()
}
