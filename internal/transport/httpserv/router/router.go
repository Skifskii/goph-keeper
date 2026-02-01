package router

import (
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/handlers/secrets"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
}

func New() *Router {
	r := chi.NewRouter()

	// middlewares
	// TODO:

	// handlers
	r.Route("/api", func(r chi.Router) {
		r.Post("/secrets", secrets.NewPost())
		r.Get("/secrets", secrets.NewGet())
	})

	return &Router{r}
}
