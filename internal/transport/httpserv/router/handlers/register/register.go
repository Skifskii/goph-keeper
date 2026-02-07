package registerhttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/repository"
)

// Registerer defines the subset of repository behavior required by the
// register HTTP handler: creating a new user and returning its ID.
type Registerer interface {
	Register(username, password string) (int, error)
}

// NewPost returns an HTTP handler for POST /register that creates a new
// user using the provided Registerer and responds with appropriate
// status codes on conflict or error.
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

// RegisterReq is the JSON payload expected by the register endpoint.
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
