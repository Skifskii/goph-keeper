package secret

import (
	"encoding/json"
	"fmt"
)

// BaseSecret contains non-sensitive metadata shared by all secret types.
// It is suitable for listings and as the base for concrete secret values.
type BaseSecret struct {
	ID         int
	UserID     int
	Metadata   string
	SecretType SecretType
}

// DecryptedSecret represents a secret with its payload in cleartext
// (not encrypted). Use this type when performing validation or before
// encrypting for storage.
type DecryptedSecret struct {
	BaseSecret
	DecPayload json.RawMessage
}

// EncryptedSecret contains a secret payload encrypted and ready for
// persistent storage. The EncPayload field holds the ciphertext.
type EncryptedSecret struct {
	BaseSecret
	EncPayload []byte
}

// NewDecryptedSecret validates payload according to the provided
// BaseSecret.SecretType and returns a DecryptedSecret on success.
// It returns a descriptive error when validation fails.
func NewDecryptedSecret(base BaseSecret, payload json.RawMessage) (DecryptedSecret, error) {
	if err := validatePayload(payload, base.SecretType); err != nil {
		return DecryptedSecret{}, fmt.Errorf("payload validation error: %w", err)
	}

	// TODO: validate payload len
	// TODO: validate metadata len

	return DecryptedSecret{
		BaseSecret: base,
		DecPayload: payload,
	}, nil
}
