package secretservice

import (
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	cryptoservice "github.com/Skifskii/goph-keeper/internal/service/secret/crypto"
)

type SecretService struct {
	repo      Repository
	encryptor Encryptor
}

type Repository interface {
	SaveSecret(payload []byte, secretType secret.SecretType, userID int, metadata string) (secretID int, err error)
}

type Encryptor interface {
	Encrypt(payload []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (payload []byte, err error)
}

func New(repo Repository, masterKey []byte) (*SecretService, error) {
	c, err := cryptoservice.New(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize crypto service: %w", err)
	}

	return &SecretService{
		repo:      repo,
		encryptor: c,
	}, nil
}

func (s *SecretService) Create(payload []byte, secretType secret.SecretType, userID int, metadata string) (id int, err error) {
	encrypted, err := s.encryptor.Encrypt(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	secretID, err := s.repo.SaveSecret(
		encrypted,
		secretType,
		userID,
		metadata,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	return secretID, nil
}
