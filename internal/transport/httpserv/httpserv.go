package httpserv

import (
	"encoding/json"
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
	CreateSecret(
		payload json.RawMessage,
		secretType string,
		metadata string,
		userID int,
	) (id int, err error)
}

type SecretGetter interface {
	GetSecret(secretID, requesterID int) (enc secret.DecryptedSecret, err error)
}

func New(log *slog.Logger, addr string, secretCreator SecretCreator, secretGetter SecretGetter) *HTTPServer {
	h := HTTPServer{
		log: log,
		server: &http.Server{
			Addr:    addr,
			Handler: router.New(log, secretCreator, secretGetter),
		},
	}

	return &h
}

func (h *HTTPServer) Run() error {
	h.log.Info("Starting HTTP server", slog.String("address", h.server.Addr))
	return h.server.ListenAndServe()
}
