package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

type Authenticator interface {
	AuthenticateWithJWT(jwtTokenString string) (userID int, err error)
}

func Auth(log *slog.Logger, authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get cookie
			cookie, _ := r.Cookie("jwt")
			if cookie == nil {
				log.Error("empty cookie")
				http.Error(w, "empty cookie", http.StatusUnauthorized)
				return
			}

			// try to get userID
			userID, err := authenticator.AuthenticateWithJWT(cookie.Value)
			if err != nil {
				log.Error("failed to authenticate with jwt", slog.Any("error", err))
				http.Error(w, "authentication failed", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
