package secretservice

import (
	"encoding/json"
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
)

type SecretService struct {
	repo      Repository
	encryptor Encryptor
}

type Repository interface {
	SaveSecret(payload []byte, secretType secret.SecretType, userID int, metadata string) (secretID int, err error)
}

type Encryptor interface {
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}

func New(repo Repository, encryptor Encryptor, masterKey []byte) (*SecretService, error) {
	return &SecretService{
		repo:      repo,
		encryptor: encryptor,
	}, nil
}

func (s *SecretService) CreateSecret(payload *secret.Payload, userID int, metadata string) (id int, err error) {
	// create domain secret
	sec, err := secret.New(payload, userID, metadata)
	if err != nil {
		return 0, fmt.Errorf("failed to initialize secret object: %w", err)
	}

	// encrypt secret payload
	plaintext, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshall payload: %w", err)
	}
	encrypted, err := s.encryptor.Encrypt(plaintext)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	// save encrypted payload to repo
	secretID, err := s.repo.SaveSecret(
		encrypted,
		sec.Payload.Type,
		sec.UserID,
		sec.Metadata,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	return secretID, nil
}
