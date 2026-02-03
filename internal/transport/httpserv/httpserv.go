package httpserv

import (
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router"
)

type HTTPServer struct {
	log    *slog.Logger
	server *http.Server
}

type SecretCreator interface {
	Create(payload []byte, secretType secret.SecretType, userID int, metadata string) (id int, err error)
}

func New(log *slog.Logger, addr string, secretCreator SecretCreator) *HTTPServer {
	h := HTTPServer{
		log: log,
		server: &http.Server{
			Addr:    addr,
			Handler: router.New(log, secretCreator),
		},
	}

	return &h
}

func (h *HTTPServer) Run() error {
	h.log.Info("Starting HTTP server", slog.String("address", h.server.Addr))
	return h.server.ListenAndServe()
}
