package secret

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrEmptyValue = errors.New("empty value")
)

type CredentialPayload struct {
	Login    string
	Password string
}

type TextPayload struct {
	Text string
}

type BinaryPayload struct {
	Data []byte
}

type CardPayload struct {
	CardNumber     string
	CardholderName string
	ExpMonth       string
	ExpYear        string
	CVC            string
}

func validatePayload(payload json.RawMessage, secretType SecretType) error {
	switch secretType {
	case TypeCredential:
		var p CredentialPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return fmt.Errorf("failed to unmarshall json (credential): %w", err)
		}

		if p.Login == "" || p.Password == "" {
			return ErrEmptyValue
		}

		return nil

	case TypeText:
		var p TextPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return fmt.Errorf("failed to unmarshall json (text): %w", err)
		}
		return nil

	case TypeBinary:
		var p BinaryPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return fmt.Errorf("failed to unmarshall json (binary): %w", err)
		}
		return nil

	case TypeCard:
		var p CardPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return fmt.Errorf("failed to unmarshall json (card): %w", err)
		}
		return nil

	default:
		return ErrUnknownSecretType
	}
}
