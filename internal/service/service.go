package service

import (
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	authservice "github.com/Skifskii/goph-keeper/internal/service/auth"
	secretservice "github.com/Skifskii/goph-keeper/internal/service/secret"
)

type Service struct {
	Secret *secretservice.SecretService
	Auth   *authservice.AuthService
}

type AuthRepository interface {
	SaveUser(username, passwordHash string) (int, error)
}

type SecretRepository interface {
	SaveSecret(enc secret.EncryptedSecret) (secretID int, err error)
	GetSecret(secretID int) (enc secret.EncryptedSecret, err error)
}

type Encryptor interface {
	Encrypt(payload []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (payload []byte, err error)
}

func New(
	authRepo AuthRepository,
	secretRepo SecretRepository,
	encryptor Encryptor,
	masterKey []byte,
) (*Service, error) {
	secretService, err := secretservice.New(secretRepo, encryptor, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret service: %w", err)
	}

	return &Service{
		Secret: secretService,
		Auth:   authservice.New(authRepo),
	}, nil
}
