package secrethttp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
)

type SecretCreator interface {
	CreateSecret(
		payload json.RawMessage,
		secretType string,
		metadata string,
		userID int,
	) (id int, err error)
}

func NewPost(log *slog.Logger, secretCreator SecretCreator) http.HandlerFunc {
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

		// get userID
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			log.Error("can't get userID")
			http.Error(w, "can't get userID", http.StatusUnauthorized)
			return
		}
		log = log.With(slog.Int("user_id", userID))

		// read the request
		var req CreateSecretReq
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			log.Error("failed to decode json body", slog.Any("error", err))
			http.Error(w, "failed to decode json body", http.StatusBadRequest)
			return
		}

		// process
		secretID, err := secretCreator.CreateSecret(
			req.Payload,
			req.SecretType,
			req.Metadata,
			userID,
		)
		if err != nil {
			log.Error("failed to create secret", slog.Any("error", err))
			http.Error(w, "failed to create secret", http.StatusInternalServerError)
			return
		}

		// compose the response
		resp := CreateSecretResp{
			ID: secretID,
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(resp)
	}
}

type CreateSecretReq struct {
	SecretType string          `json:"secret_type"`
	Metadata   string          `json:"metadata"`
	Payload    json.RawMessage `json:"payload"`
}

type CreateSecretResp struct {
	ID int
}
