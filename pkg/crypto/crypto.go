package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// ErrCiphertextTooShort is returned when provided ciphertext is smaller
// than the expected nonce size and therefore invalid for decryption.
var ErrCiphertextTooShort = errors.New("ciphertext too short")

// ErrDecryptFailed is returned when authenticated decryption fails.
var ErrDecryptFailed = errors.New("failed to decrypt")

const (
	nonceSize     = 12
	masterKeySize = 32
)

// Crypto provides authenticated encryption and decryption using AES-GCM.
// It is configured with a fixed-size master key and exposes simple
// Encrypt/Decrypt methods that operate on byte slices.
type Crypto struct {
	aead cipher.AEAD
}

// New creates a new Crypto instance using the provided masterKey.
// The masterKey must be exactly 32 bytes long; otherwise an error is returned.
func New(masterKey []byte) (*Crypto, error) {
	if len(masterKey) != masterKeySize {
		return nil, fmt.Errorf("invalid master key length: got %d, expected %d", len(masterKey), masterKeySize)
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	return &Crypto{aead: aead}, nil
}

// Encrypt encrypts and authenticates the given payload and returns
// the concatenation of nonce and ciphertext suitable for storage or transport.
func (c *Crypto) Encrypt(payload []byte) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := c.aead.Seal(
		nil,
		nonce,
		payload,
		nil,
	)

	result := make([]byte, 0, nonceSize+len(ciphertext))
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt verifies and decrypts data previously produced by Encrypt.
// It expects the input to contain the nonce followed by the ciphertext.
func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	if len(data) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	payload, err := c.aead.Open(
		nil,
		nonce,
		ciphertext,
		nil,
	)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return payload, nil
}
