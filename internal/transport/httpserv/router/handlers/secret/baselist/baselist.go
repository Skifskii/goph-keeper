package baselisthttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
)

// BaseSecretsLister defines the contract for listing non-sensitive
// secret metadata for a user with pagination support.
type BaseSecretsLister interface {
	GetBaseSecretsList(userID, limit, offset int) ([]secret.BaseSecret, error)
}

// NewGet returns an HTTP handler for GET /secret/baselist that returns a
// paginated list of base secret metadata for the authenticated user.
func NewGet(log *slog.Logger, lister BaseSecretsLister) http.HandlerFunc {
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

		// parse query params
		q := r.URL.Query()

		limit := 50 // default
		if l := q.Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		offset := 0 // default
		if o := q.Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		// call service
		secrets, err := lister.GetBaseSecretsList(userID, limit, offset)
		if err != nil {
			log.Error("failed to list secrets", slog.Any("error", err))
			http.Error(w, "failed to list secrets", http.StatusInternalServerError)
			return
		}

		// prepare response
		resp := make([]BaseSecretResponse, len(secrets))
		for i, s := range secrets {
			resp[i] = BaseSecretResponse{
				ID:         s.ID,
				Metadata:   s.Metadata,
				SecretType: string(s.SecretType),
			}
		}

		// serialize
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(resp); err != nil {
			log.Error("failed to encode response", slog.Any("error", err))
		}
	}
}

// BaseSecretResponse is the JSON representation of the non-sensitive
// secret metadata returned in a listing response.
type BaseSecretResponse struct {
	ID         int    `json:"id"`
	SecretType string `json:"secret_type"`
	Metadata   string `json:"metadata"`
}
