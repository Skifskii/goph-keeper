package cryptoservice

import (
	"bytes"
	"testing"
)

func TestCryptoService_EncryptDecrypt_E2E(t *testing.T) {
	masterKey := []byte("super-secret-master-key-32bytes!") // пример
	payload := []byte("hello, encrypted world")

	service, err := New(masterKey)
	if err != nil {
		t.Fatalf("failed to create CryptoService: %v", err)
	}

	ciphertext, err := service.Encrypt(payload)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, payload) {
		t.Fatal("ciphertext must differ from plaintext")
	}

	decrypted, err := service.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, payload) {
		t.Fatalf(
			"decrypted payload mismatch:\nexpected: %q\ngot:      %q",
			payload,
			decrypted,
		)
	}
}
