package router

import (
	"encoding/json"
	"log/slog"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	secrethttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
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

func New(log *slog.Logger, secretCreator SecretCreator, secretGetter SecretGetter) *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/secret", secrethttp.NewPost(log, secretCreator))
		r.Get("/secret/{id}", secrethttp.NewGet(log, secretGetter))
	})

	return &Router{r}
}
