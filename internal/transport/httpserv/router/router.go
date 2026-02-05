package router

import (
	"log/slog"

	loginhttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/login"
	registerhttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/register"
	secrethttp "github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secret"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
}

type Auther interface {
	registerhttp.Registerer
	loginhttp.Loginer
	middleware.Authenticator
}

func New(
	log *slog.Logger,
	auther Auther,
	secretCreator secrethttp.SecretCreator,
	secretGetter secrethttp.SecretGetter,
) *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", registerhttp.NewPost(log, auther))
		r.Post("/login", loginhttp.NewPost(log, auther))

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(log, auther))

			r.Route("/secret", func(r chi.Router) {
				r.Post("/", secrethttp.NewPost(log, secretCreator))
				r.Get("/{id}", secrethttp.NewGet(log, secretGetter))
			})
		})
	})

	return &Router{r}
}
