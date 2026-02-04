package router

import (
	"encoding/json"
	"log/slog"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	registerhttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/register"
	secrethttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
}

type Auther interface {
	Register(username, password string) (int, error)
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

func New(
	log *slog.Logger,
	auther Auther,
	secretCreator SecretCreator,
	secretGetter SecretGetter,
) *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", registerhttp.NewPost(log, auther))

		r.Route("/secret", func(r chi.Router) {
			r.Post("/", secrethttp.NewPost(log, secretCreator))
			r.Get("/{id}", secrethttp.NewGet(log, secretGetter))
		})
	})

	return &Router{r}
}
