package secrethttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/repository"
	"github.com/go-chi/chi/v5"
)

type SecretGetter interface {
	GetSecret(secretID, requesterID int) (enc secret.DecryptedSecret, err error)
}

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

		userID := 1 // TODO: add middleware Auth

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

type GetSecretResponse struct {
	SecretType string          `json:"secret_type"`
	Metadata   string          `json:"metadata"`
	Payload    json.RawMessage `json:"payload"`
}
