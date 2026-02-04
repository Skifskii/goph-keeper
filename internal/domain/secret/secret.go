package secret

import (
	"encoding/json"
	"fmt"
)

type BaseSecret struct {
	UserID     int
	Metadata   string
	SecretType SecretType
}

type DecryptedSecret struct {
	BaseSecret
	DecPayload json.RawMessage
}

type EncryptedSecret struct {
	BaseSecret
	EncPayload []byte
}

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
