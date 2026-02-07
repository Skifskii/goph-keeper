package secrethttp

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Skifskii/goph-keeper/internal/repository"
	secretservice "github.com/Skifskii/goph-keeper/internal/service/secret"
	"github.com/Skifskii/goph-keeper/internal/transport/httpserv/router/middleware"
	"github.com/go-chi/chi/v5"
)

type SecretDeleter interface {
	DeleteSecret(secretID, userID int) error
}

func NewDelete(log *slog.Logger, secretDeleter SecretDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
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
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		log = log.With(slog.Int("user_id", userID))

		log.Info("new request")

		// get secretID
		secretIDParam := chi.URLParam(r, "id")
		if secretIDParam == "" {
			http.Error(w, "missing secret id", http.StatusBadRequest)
			return
		}

		secretID, err := strconv.Atoi(secretIDParam)
		if err != nil {
			http.Error(w, "invalid secret id", http.StatusBadRequest)
			return
		}

		// process
		err = secretDeleter.DeleteSecret(secretID, userID)
		if err != nil {
			switch {
			case errors.Is(err, secretservice.ErrSecretAccessDenied):
				http.Error(w, "access denied", http.StatusForbidden)
				return
			case errors.Is(err, repository.ErrSecretNotFound):
				http.Error(w, "secret not found", http.StatusNotFound)
				return
			default:
				log.Error("failed to delete secret", slog.Any("error", err))
				http.Error(w, "failed to delete secret", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
