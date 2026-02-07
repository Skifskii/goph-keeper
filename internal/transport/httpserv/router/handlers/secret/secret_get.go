package secrethttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/repository"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
	"github.com/go-chi/chi/v5"
)

// SecretGetter describes the service behavior required by the GET
// /secret/{id} handler: fetching and returning a decrypted secret
// for a specific requester.
type SecretGetter interface {
	GetSecret(secretID, requesterID int) (enc secret.DecryptedSecret, err error)
}

// NewGet returns an HTTP handler for GET /secret/{id} that verifies
// ownership and returns the decrypted secret payload as JSON.
func NewGet(log *slog.Logger, secretGetter SecretGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log = log.With(
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
		)

		// get userID
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			log.Error("can't get userID")
			http.Error(w, "can't get userID", http.StatusUnauthorized)
			return
		}
		log = log.With(slog.Int("user_id", userID))

		log.Info("new request")

		// read params
		secretIDParam := chi.URLParam(r, "id")
		if secretIDParam == "" {
			log.Error("failed to get secret id from params")
			http.Error(w, "failed to get secret id from params", http.StatusBadRequest)
			return
		}
		secretID, err := strconv.Atoi(secretIDParam)
		if err != nil {
			log.Error("failed to convert secretID to int", slog.Any("error", err))
			http.Error(w, "failed to convert secretID to int", http.StatusBadRequest)
			return
		}

		// process
		sec, err := secretGetter.GetSecret(secretID, userID)
		if err != nil {
			if errors.Is(err, repository.ErrSecretNotFound) {
				log.Error("failed to get secret", slog.Any("error", err))
				http.Error(w, "secret not found", http.StatusNotFound)
				return
			}
			log.Error("failed to get secret", slog.Any("error", err))
			http.Error(w, "failed to get secret", http.StatusInternalServerError)
			return
		}

		// compose the response
		resp := GetSecretResponse{
			SecretType: string(sec.SecretType),
			Metadata:   sec.Metadata,
			Payload:    sec.DecPayload,
		}

		// serialize the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(resp)
	}
}

// GetSecretResponse is the JSON response body returned by GET /secret/{id}.
type GetSecretResponse struct {
	SecretType string          `json:"secret_type"`
	Metadata   string          `json:"metadata"`
	Payload    json.RawMessage `json:"payload"`
}
