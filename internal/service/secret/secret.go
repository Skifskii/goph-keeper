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
)

type SecretService struct {
	repo      Repository
	encryptor Encryptor
}

type Repository interface {
	SaveSecret(enc secret.EncryptedSecret) (secretID int, err error)
	GetSecret(secretID int) (enc secret.EncryptedSecret)
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
		return 0, fmt.Errorf("failed to initialize decrypted secret object: %w", err)
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
	return secret.DecryptedSecret{}, nil // TODO: implement
	// 	// get secret from repo
	// 	ciphertext, secretType, ownerID, metadata, err := s.repo.GetSecret(secretID)
	// 	if err != nil {
	// 		return secret.Secret{}, fmt.Errorf("failed to get secret from repo: %w", err)
	// 	}

	// 	// check if the requester has access to the secret
	// 	if requesterID != ownerID {
	// 		return secret.Secret{}, ErrSecretAccessDenied
	// 	}

	// 	// decrypt secret payload
	// 	var payload secret.Payload
	// 	plaintext, err := s.encryptor.Decrypt(ciphertext)
	// 	if err != nil {
	// 		return secret.Secret{}, fmt.Errorf("failed to decrypt secret payload: %w", err)
	// 	}
	// 	if err := json.Unmarshal(plaintext, &payload); err != nil {
	// 		return secret.Secret{}, fmt.Errorf("failed to unmarshall plaintext: %w", err)
	// 	}

	// 	// validate secret type
	// 	if secretType != payload.Type {
	// 		return secret.Secret{}, ErrDifferentTypes
	// 	}

	//	return secret.Secret{
	//		Payload:  &payload,
	//		UserID:   ownerID,
	//		Metadata: metadata,
	//	}, nil
}
