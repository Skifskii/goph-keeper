package httpserv

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router"
	secrethttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret"
	baselisthttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret/baselist"
)

type HTTPServer struct {
	log    *slog.Logger
	server *http.Server
}

func New(
	log *slog.Logger, addr string,
	auther router.Auther,
	secretCreator secrethttp.SecretCreator,
	secretGetter secrethttp.SecretGetter,
	baseSecretLister baselisthttp.BaseSecretsLister,
) *HTTPServer {
	h := HTTPServer{
		log: log,
		server: &http.Server{
			Addr:    addr,
			Handler: router.New(log, auther, secretCreator, secretGetter, baseSecretLister),
		},
	}

	return &h
}

func (h *HTTPServer) MustRun() {
	err := h.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (h *HTTPServer) Run() error {
	h.log.Info("Starting HTTP server", slog.String("address", h.server.Addr))
	return h.server.ListenAndServe()
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
