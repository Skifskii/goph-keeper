package secret

import "errors"

// ErrUnknownSecretType is returned when an unrecognized secret type
// string is encountered.
var ErrUnknownSecretType = errors.New("unknown secret type")

// SecretType represents the category of a secret payload.
type SecretType string

// Supported SecretType values.
const (
	// TypeCredential denotes a username/password credential payload.
	TypeCredential SecretType = "credential"
	// TypeText denotes a plain text payload.
	TypeText SecretType = "text"
	// TypeBinary denotes an opaque binary payload.
	TypeBinary SecretType = "binary"
	// TypeCard denotes a payment card-like payload.
	TypeCard SecretType = "card"
)

// NewSecretTypeFromString converts a string to a SecretType, returning
// ErrUnknownSecretType for unsupported values.
func NewSecretTypeFromString(s string) (SecretType, error) {
	switch s {
	case string(TypeCredential):
		return TypeCredential, nil
	case string(TypeText):
		return TypeText, nil
	case string(TypeBinary):
		return TypeBinary, nil
	case string(TypeCard):
		return TypeCard, nil
	default:
		return "", ErrUnknownSecretType
	}
}

// ToInt maps a SecretType to a stable integer identifier used by the
// persistence layer. It returns ErrUnknownSecretType for unsupported types.
func (st SecretType) ToInt() (int, error) {
	switch st {
	case TypeCredential:
		return 1, nil
	case TypeText:
		return 2, nil
	case TypeBinary:
		return 3, nil
	case TypeCard:
		return 4, nil
	default:
		return 0, ErrUnknownSecretType
	}
}
