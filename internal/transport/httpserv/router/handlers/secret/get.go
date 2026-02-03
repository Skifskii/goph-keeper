package secrethttp

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SecretGetter interface {
	GetSecret(secretID, userID int) ()
}

func NewGet(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		secretID := chi.URLParam(r, "id")
		if secretID == "" {
			log.Error("failed to get secret id from params")
			http.Error(w, "failed to get secret id from params", http.StatusBadRequest)
			return
		}



	}
}
