package router

import (
	"log/slog"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	secrethttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
}

type SecretCreator interface {
	Create(payload []byte, secretType secret.SecretType, userID int, metadata string) (id int, err error)
}

func New(log *slog.Logger, secretCreator SecretCreator) *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/secrets", secrethttp.NewPost(log, secretCreator))
		r.Get("/secrets", secrethttp.NewGet())
	})

	return &Router{r}
}
