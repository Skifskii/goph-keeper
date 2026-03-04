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

// HTTPServer wraps an http.Server and its dependencies for serving the
// application's HTTP API. It exposes lifecycle helpers for running and
// stopping the server.
type HTTPServer struct {
	log    *slog.Logger
	server *http.Server
}

// New constructs an HTTPServer configured with routes produced by the
// router package. Pass handler adapters and an Auther implementation
// required by the router.
func New(
	log *slog.Logger, addr string,
	auther router.Auther,
	secretCreator secrethttp.SecretCreator,
	secretGetter secrethttp.SecretGetter,
	secretUpdater secrethttp.SecretUpdater,
	secretDeleter secrethttp.SecretDeleter,
	baseSecretLister baselisthttp.BaseSecretsLister,
) *HTTPServer {
	h := HTTPServer{
		log: log,
		server: &http.Server{
			Addr:    addr,
			Handler: router.New(log, auther, secretCreator, secretGetter, secretUpdater, secretDeleter, baseSecretLister),
		},
	}

	return &h
}

// MustRun runs the server and panics on non-graceful errors. It is a
// convenience wrapper used by the application's top-level runner.
func (h *HTTPServer) MustRun() {
	err := h.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

// Run starts listening and serving HTTP requests. It returns any error
// produced by ListenAndServe.
func (h *HTTPServer) Run() error {
	h.log.Info("Starting HTTP server", slog.String("address", h.server.Addr))
	return h.server.ListenAndServe()
}

// Stop gracefully shuts down the server using the provided context.
func (h *HTTPServer) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
