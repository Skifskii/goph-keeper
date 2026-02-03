package secret

import (
	"errors"
)

var (
	ErrUnknownSecretType = errors.New("unknown secret type")
)

type Secret struct {
	ID      int
	OwnerID int
	Type    SecretType
	Payload []byte
}

type SecretType string

const (
	TypeCredential SecretType = "credential"
	TypeText       SecretType = "text"
	TypeBinary     SecretType = "binary"
	TypeCard       SecretType = "card"
)

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
