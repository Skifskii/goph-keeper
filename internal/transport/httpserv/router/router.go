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
	CreateSecret(payload *secret.Payload, userID int, metadata string) (id int, err error)
}

func New(log *slog.Logger, secretCreator SecretCreator) *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/secret", secrethttp.NewPost(log, secretCreator))
		r.Get("/secret{id}", secrethttp.NewGet(log))
	})

	return &Router{r}
}
