package httpserv

import (
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router"
)

type HTTPServer struct {
	log    *slog.Logger
	server *http.Server
}

func New(log *slog.Logger, addr string) *HTTPServer {
	h := HTTPServer{
		log: log,
		server: &http.Server{
			Addr:    addr,
			Handler: router.New(),
		},
	}

	return &h
}

func (h *HTTPServer) Run() error {
	h.log.Info("Starting HTTP server", slog.String("address", h.server.Addr))
	return h.server.ListenAndServe()
}
