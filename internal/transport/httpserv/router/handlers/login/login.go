package loginhttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	authservice "github.com/Skifskii/goph-keeper/internal/service/auth"
)

// Loginer defines the subset of the auth service required by the
// login HTTP handler: authenticating credentials and returning a JWT.
type Loginer interface {
	Login(username, password string) (string, error)
}

// NewPost returns an HTTP handler for POST /login that validates user
// credentials via the provided Loginer and sets a JWT cookie on success.
func NewPost(log *slog.Logger, loginer Loginer) http.HandlerFunc {
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
		var req LoginReq
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			log.Error("failed to decode json body", slog.Any("error", err))
			http.Error(w, "failed to decode json body", http.StatusBadRequest)
			return
		}

		// process
		jwtToken, err := loginer.Login(req.Username, req.Password)
		if err != nil {
			if errors.Is(err, authservice.ErrInvalidCredentials) {
				log.Error("failed to login", slog.Any("error", err))
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			log.Error("failed to login", slog.Any("error", err))
			http.Error(w, "failed to login", http.StatusInternalServerError)
			return
		}

		// compose the response
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: jwtToken,
		})
		w.WriteHeader(http.StatusOK)
	}
}

// LoginReq represents the JSON payload expected by the login endpoint.
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
