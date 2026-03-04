package user

// User represents a persisted user account with credential information.
// Fields hold identifying data and the stored password hash; methods
// operating on users live elsewhere in the codebase.
type User struct {
	ID           int
	Username     string
	PasswordHash []byte
}
