package registerhttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/repository"
)

type Registerer interface {
	Register(username, password string) (int, error)
}

func NewPost(log *slog.Logger, registerer Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log = log.With(
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
		)

		// read the request
		var req RegisterReq
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			log.Error("failed to decode json body", slog.Any("error", err))
			http.Error(w, "failed to decode json body", http.StatusBadRequest)
			return
		}

		// process
		_, err := registerer.Register(req.Username, req.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUsernameTaken) {
				log.Error("failed to register user", slog.Any("error", err))
				http.Error(w, "username already taken", http.StatusConflict)
				return
			}

			log.Error("failed to register user", slog.Any("error", err))
			http.Error(w, "failed to register user", http.StatusInternalServerError)
			return
		}

		log.Info("new user successfully registered!")

		// compose the response
		w.WriteHeader(http.StatusOK)
	}
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
