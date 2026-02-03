package service

import (
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	secretservice "github.com/Skifskii/goph-keeper/internal/service/secret"
)

type Service struct {
	Secret *secretservice.SecretService
}

type SecretSaver interface {
	SaveSecret(payload []byte, secretType secret.SecretType, userID int, metadata string) (secretID int, err error)
}

func New(secretSaver SecretSaver, masterKey []byte) (*Service, error) {
	secretService, err := secretservice.New(secretSaver, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret service: %w", err)
	}

	return &Service{
		Secret: secretService,
	}, nil
}
