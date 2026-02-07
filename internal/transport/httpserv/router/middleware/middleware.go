package middleware

// ContextKey is the type used for request context keys in middleware to
// avoid collisions with other context users.
type ContextKey string

// UserIDKey is the context key under which authenticated user ID is
// stored by the authentication middleware.
const UserIDKey ContextKey = "user_id"
