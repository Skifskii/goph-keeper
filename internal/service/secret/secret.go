package secretservice

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
)

var (
	ErrSecretAccessDenied = errors.New("secret access denied")
	ErrDifferentTypes     = errors.New("secret type from repo don't match payload type")
	ErrRequestValidation  = errors.New("validation error")
)

type SecretService struct {
	repo      Repository
	encryptor Encryptor
}

type Repository interface {
	SaveSecret(enc secret.EncryptedSecret) (secretID int, err error)
	GetSecret(secretID int) (enc secret.EncryptedSecret, err error)
	GetUserSecrets(userID, limit, offset int) ([]secret.BaseSecret, error)
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

func (s *SecretService) CreateSecret(
	payload json.RawMessage,
	secretTypeString string,
	metadata string,
	userID int,
) (id int, err error) {
	// create domain secret
	secretType, err := secret.NewSecretTypeFromString(secretTypeString)
	if err != nil {
		return 0, fmt.Errorf("secret type validation error: %w", err)
	}
	decSecret, err := secret.NewDecryptedSecret(
		secret.BaseSecret{
			UserID:     userID,
			Metadata:   metadata,
			SecretType: secretType,
		},
		payload,
	)
	if err != nil {
		return 0, ErrRequestValidation
	}

	// encrypt secret
	encSecret, err := s.encryptSecret(decSecret)
	if err != nil {
		return 0, err
	}

	// save encrypted payload to repo
	secretID, err := s.repo.SaveSecret(encSecret)
	if err != nil {
		return 0, fmt.Errorf("failed to save secret to repo: %w", err)
	}

	return secretID, nil
}

func (s *SecretService) encryptSecret(dec secret.DecryptedSecret) (secret.EncryptedSecret, error) {
	encPayload, err := s.encryptor.Encrypt(dec.DecPayload)
	if err != nil {
		return secret.EncryptedSecret{}, fmt.Errorf("failed to encrypt payload: %w", err)
	}
	return secret.EncryptedSecret{
		BaseSecret: dec.BaseSecret,
		EncPayload: encPayload,
	}, nil
}

func (s *SecretService) GetSecret(secretID, requesterID int) (secret.DecryptedSecret, error) {
	// get secret from repo
	encSecret, err := s.repo.GetSecret(secretID)
	if err != nil {
		return secret.DecryptedSecret{}, fmt.Errorf("failed to get secret from repo: %w", err)
	}

	// check if the requester has access to the secret
	if requesterID != encSecret.UserID {
		return secret.DecryptedSecret{}, ErrSecretAccessDenied
	}

	// decrypt secret payload
	decSecret, err := s.decryptSecret(encSecret)
	if err != nil {
		return secret.DecryptedSecret{}, fmt.Errorf("failed to decrypt secret payload: %w", err)
	}

	return decSecret, nil
}

func (s *SecretService) GetBaseSecretsList(userID, limit, offset int) ([]secret.BaseSecret, error) {
	secrets, err := s.repo.GetUserSecrets(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user secrets from repo: %w", err)
	}

	return secrets, nil
}

func (s *SecretService) decryptSecret(dec secret.EncryptedSecret) (secret.DecryptedSecret, error) {
	decPayload, err := s.encryptor.Decrypt(dec.EncPayload)
	if err != nil {
		return secret.DecryptedSecret{}, fmt.Errorf("failed to decrypt payload: %w", err)
	}
	return secret.DecryptedSecret{
		BaseSecret: dec.BaseSecret,
		DecPayload: decPayload,
	}, nil
}
