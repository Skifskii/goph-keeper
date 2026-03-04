package service

import (
	"fmt"
	"time"

	authservice "github.com/Skifskii/goph-keeper/internal/service/auth"
	secretservice "github.com/Skifskii/goph-keeper/internal/service/secret"
)

// Service aggregates application services used by the server.
// It exposes initialized service instances such as Secret and Auth.
type Service struct {
	Secret *secretservice.SecretService
	Auth   *authservice.AuthService
}

// New constructs and initializes the main Service bundle.
// It wires the provided repositories, encryptor and keys into concrete
// service implementations and returns an error on initialization failure.
func New(
	authRepo authservice.Repository,
	secretRepo secretservice.Repository,
	encryptor secretservice.Encryptor,
	masterKey []byte,
	secretKey string,
	jwtTokenTTL time.Duration,
) (*Service, error) {
	secretService, err := secretservice.New(secretRepo, encryptor, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret service: %w", err)
	}

	return &Service{
		Secret: secretService,
		Auth:   authservice.New(authRepo, secretKey, jwtTokenTTL),
	}, nil
}
