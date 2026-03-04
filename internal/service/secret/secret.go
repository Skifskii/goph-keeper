package secretservice

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
)

// ErrSecretAccessDenied indicates the current requester is not authorized
// to access the requested secret.
var ErrSecretAccessDenied = errors.New("secret access denied")

// ErrDifferentTypes is returned when a secret's stored type does not
// match the provided payload type.
var ErrDifferentTypes = errors.New("secret type from repo don't match payload type")

// ErrRequestValidation indicates the incoming request payload failed
// domain validation.
var ErrRequestValidation = errors.New("validation error")

// SecretService provides business operations for creating, retrieving,
// updating and deleting secrets. It relies on a repository for storage
// and an encryptor for payload protection.
type SecretService struct {
	repo      Repository
	encryptor Encryptor
}

// Repository defines storage operations required by SecretService.
// Implementations are expected to persist encrypted secrets and expose
// retrieval and mutation methods used by the service.
type Repository interface {
	SaveSecret(enc secret.EncryptedSecret) (secretID int, err error)
	GetSecret(secretID int) (enc secret.EncryptedSecret, err error)
	GetUserSecrets(userID, limit, offset int) ([]secret.BaseSecret, error)
	UpdateSecret(enc secret.EncryptedSecret) error
	DeleteSecret(secretID int) error
}

// Encryptor abstracts authenticated encryption operations used by the
// service to protect secret payloads prior to persistence.
type Encryptor interface {
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}

// New constructs a SecretService with the given repository and encryptor.
// The masterKey parameter is accepted for symmetry with constructors and
// may be required by some implementations of Encryptor; it is not used
// directly by this function.
func New(repo Repository, encryptor Encryptor, masterKey []byte) (*SecretService, error) {
	return &SecretService{
		repo:      repo,
		encryptor: encryptor,
	}, nil
}

// CreateSecret validates and stores a new secret for the given user.
// It returns the created secret ID or an error if validation, encryption,
// or repository persistence fail.
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

// GetSecret returns a decrypted secret identified by secretID if the
// requesterID has access to it. It returns ErrSecretAccessDenied when
// the requester is not the owner.
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

// GetBaseSecretsList returns a paginated list of base secret metadata
// for the specified user. The returned slice contains non-sensitive
// fields appropriate for listing views.
func (s *SecretService) GetBaseSecretsList(userID, limit, offset int) ([]secret.BaseSecret, error) {
	secrets, err := s.repo.GetUserSecrets(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user secrets from repo: %w", err)
	}

	return secrets, nil
}

// UpdateSecret validates and updates an existing secret's payload and
// metadata. It enforces ownership and returns the secret ID on success.
func (s *SecretService) UpdateSecret(
	payload json.RawMessage,
	metadata string,
	secretID, userID int,
) (int, error) {
	// get secret from repo
	encSecret, err := s.repo.GetSecret(secretID)
	if err != nil {
		return 0, fmt.Errorf("failed to get secret from repo: %w", err)
	}

	// check access
	if encSecret.UserID != userID {
		return 0, ErrSecretAccessDenied
	}

	// create new decrypted secret
	decSecret, err := secret.NewDecryptedSecret(
		secret.BaseSecret{
			ID:         encSecret.ID,
			UserID:     encSecret.UserID,
			SecretType: encSecret.SecretType, // тип менять нельзя
			Metadata:   metadata,
		},
		payload,
	)
	if err != nil {
		return 0, ErrRequestValidation
	}

	// encrypt payload
	updatedEncSecret, err := s.encryptSecret(decSecret)
	if err != nil {
		return 0, err
	}

	// update in repo
	if err := s.repo.UpdateSecret(updatedEncSecret); err != nil {
		return 0, fmt.Errorf("failed to update secret in repo: %w", err)
	}

	return secretID, nil
}

// DeleteSecret removes the specified secret if the provided userID
// is the owner. It returns ErrSecretAccessDenied when the user is not
// authorized to delete the secret.
func (s *SecretService) DeleteSecret(secretID, userID int) error {
	// get secret from repo
	encSecret, err := s.repo.GetSecret(secretID)
	if err != nil {
		return fmt.Errorf("failed to get secret from repo: %w", err)
	}

	// check access
	if encSecret.UserID != userID {
		return ErrSecretAccessDenied
	}

	// delete
	if err := s.repo.DeleteSecret(secretID); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}
