package secrethttp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
)

type SecretGetter interface {
	GetSecret(secretID, requesterID int) (enc secret.DecryptedSecret, err error)
}

func NewGet(log *slog.Logger, secretGetter SecretGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 	if r.Method != http.MethodGet {
		// 		w.WriteHeader(http.StatusMethodNotAllowed)
		// 		return
		// 	}

		// 	log = log.With(
		// 		slog.String("uri", r.RequestURI),
		// 		slog.String("method", r.Method),
		// 	)

		// 	userID := 1 // TODO: add middleware Auth

		// 	// read params
		// 	secretIDParam := chi.URLParam(r, "id")
		// 	if secretIDParam == "" {
		// 		log.Error("failed to get secret id from params")
		// 		http.Error(w, "failed to get secret id from params", http.StatusBadRequest)
		// 		return
		// 	}
		// 	secretID, err := strconv.Atoi(secretIDParam)
		// 	if err != nil {
		// 		log.Error("failed to convert secretID to int", slog.Any("error", err))
		// 		http.Error(w, "failed to convert secretID to int", http.StatusBadRequest)
		// 		return
		// 	}

		// 	// process
		// 	sec, err := secretGetter.GetSecret(secretID, userID)
		// 	if err != nil {
		// 		// TODO: process different errors
		// 		log.Error("failed to get secret", slog.Any("error", err))
		// 		http.Error(w, "failed to get secret", http.StatusBadRequest)
		// 		return
		// 	}

		// 	// compose the response
		// 	payladJSON, err := sec.Payload.GetJSON()
		// 	if err != nil {
		// 		log.Error("failed to build json from payload", slog.Any("error", err))
		// 		http.Error(w, "failed to build json from payload", http.StatusBadRequest)
		// 		return
		// 	}
		// 	resp := GetSecretResponse{
		// 		SecretType: string(sec.Payload.Type),
		// 		Metadata:   sec.Metadata,
		// 		Payload:    payladJSON,
		// 	}

		// 	// serialize the response
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusOK)
		// 	encoder := json.NewEncoder(w)
		// 	encoder.Encode(resp)
	}
}

type GetSecretResponse struct {
	SecretType string          `json:"secret_type"`
	Metadata   string          `json:"metadata"`
	Payload    json.RawMessage `json:"payload"`
}
