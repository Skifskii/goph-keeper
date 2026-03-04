package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

// Authenticator defines the behavior required by the authentication
// middleware: validating a JWT and returning the associated user ID.
type Authenticator interface {
	AuthorizeWithJWT(jwtTokenString string) (userID int, err error)
}

// Auth returns a middleware that extracts a JWT from the request cookie,
// validates it via the provided Authenticator and stores the resulting
// user ID in the request context under UserIDKey.
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
			userID, err := authenticator.AuthorizeWithJWT(cookie.Value)
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
