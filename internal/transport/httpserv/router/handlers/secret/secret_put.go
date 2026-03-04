package secrethttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	secretservice "github.com/Skifskii/goph-keeper/internal/service/secret"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
)

// SecretUpdater defines the contract required by the PUT /secret/{id}
// handler to update an existing secret's payload and metadata.
type SecretUpdater interface {
	UpdateSecret(
		payload json.RawMessage,
		metadata string,
		secretID, userID int,
	) (id int, err error)
}

// NewPut returns an HTTP handler for PUT /secret/{id} which validates
// ownership and applies an update via SecretUpdater.
func NewPut(log *slog.Logger, secretUpdater SecretUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPut {
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

		// get secretID from path
		secretID, err := getSecretIDFromRequest(r)
		if err != nil {
			log.Error("invalid secret id", slog.Any("error", err))
			http.Error(w, "invalid secret id", http.StatusBadRequest)
			return
		}
		log = log.With(slog.Int("secret_id", secretID))

		log.Info("update secret request")

		// read request body
		var req UpdateSecretReq
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			log.Error("failed to decode json body", slog.Any("error", err))
			http.Error(w, "failed to decode json body", http.StatusBadRequest)
			return
		}

		// process
		id, err := secretUpdater.UpdateSecret(
			req.Payload,
			req.Metadata,
			secretID,
			userID,
		)
		if err != nil {
			switch {
			case errors.Is(err, secretservice.ErrRequestValidation):
				log.Error("validation error", slog.Any("error", err))
				http.Error(w, "validation error", http.StatusBadRequest)
				return

			case errors.Is(err, secretservice.ErrSecretAccessDenied):
				log.Error("access denied", slog.Any("error", err))
				http.Error(w, "access denied", http.StatusForbidden)
				return

			default:
				log.Error("failed to update secret", slog.Any("error", err))
				http.Error(w, "failed to update secret", http.StatusInternalServerError)
				return
			}
		}

		// response
		resp := UpdateSecretResp{ID: id}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func getSecretIDFromRequest(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("missing secret id")
	}
	return strconv.Atoi(idStr)
}

// UpdateSecretReq is the JSON payload expected by the update endpoint.
type UpdateSecretReq struct {
	Metadata string          `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

// UpdateSecretResp is returned after a successful update and contains
// the id of the updated secret.
type UpdateSecretResp struct {
	ID int `json:"id"`
}
