package model

type Secret struct {
	ID      int
	OwnerID int
	Type    SecretType
}

type SecretType string

const (
	TypeCredential SecretType = "credential"
	TypeText       SecretType = "text"
	TypeBinary     SecretType = "binary"
	TypeCard       SecretType = "card"
)
