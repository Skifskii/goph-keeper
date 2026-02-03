package cryptoservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var (
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	ErrDecryptFailed      = errors.New("failed to decrypt")
)

const (
	nonceSize     = 12
	masterKeySize = 32
)

type CryptoService struct {
	aead cipher.AEAD
}

func New(masterKey []byte) (*CryptoService, error) {
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

	return &CryptoService{aead: aead}, nil
}

func (c *CryptoService) Encrypt(payload []byte) ([]byte, error) {
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

func (c *CryptoService) Decrypt(data []byte) ([]byte, error) {
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
