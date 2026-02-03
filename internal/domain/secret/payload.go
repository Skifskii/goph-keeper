package secret

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	Type       SecretType
	Credential *CredentialPayload
	Card       *CardPayload
	Text       *TextPayload
	Binary     *BinaryPayload
}

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

func NewPayloadFromJSON(jsonPayload json.RawMessage, secretType SecretType) (*Payload, error) {
	switch secretType {
	case TypeCredential:
		var p CredentialPayload
		if err := json.Unmarshal(jsonPayload, &p); err != nil {
			return nil, fmt.Errorf("failed to unmarshall json (credential): %w", err)
		}
		return &Payload{
			Type:       TypeCredential,
			Credential: &p,
		}, nil

	case TypeText:
		var p TextPayload
		if err := json.Unmarshal(jsonPayload, &p); err != nil {
			return nil, fmt.Errorf("failed to unmarshall json (text): %w", err)
		}
		return &Payload{
			Type: TypeText,
			Text: &p,
		}, nil

	case TypeBinary:
		var p BinaryPayload
		if err := json.Unmarshal(jsonPayload, &p); err != nil {
			return nil, fmt.Errorf("failed to unmarshall json (binary): %w", err)
		}
		return &Payload{
			Type:   TypeBinary,
			Binary: &p,
		}, nil

	case TypeCard:
		var p CardPayload
		if err := json.Unmarshal(jsonPayload, &p); err != nil {
			return nil, fmt.Errorf("failed to unmarshall json (card): %w", err)
		}
		return &Payload{
			Type: TypeCard,
			Card: &p,
		}, nil

	default:
		return nil, ErrUnknownSecretType
	}
}
