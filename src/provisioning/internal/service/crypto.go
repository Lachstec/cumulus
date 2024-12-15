package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
)

// CryptoService provides utility functions to handle SSH keys
// that are used in order to execute commands on the game servers.
// This service should not be directly exposed to end user requests.
type CryptoService struct {
	// encryptionKey is used to encrypt SSH keys in order to store them in the database securely.
	encryptionKey []byte
}

// NewCryptoService creates a new CryptoService with the specified private key
// as encryption key. A safe key of appropriate length should be used.
func NewCryptoService(encryptionKey []byte) *CryptoService {
	return &CryptoService{encryptionKey: encryptionKey}
}

// EncryptPrivateKey encrypts the given ed25519 private key with the secret key stored in the
// CryptoService. Returns the resulting, base64-encoded ciphertext or an error.
func (c *CryptoService) EncryptPrivateKey(key ed25519.PrivateKey) (string, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, key.Seed(), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey takes in a base64 encoded private ed25519 key and
// tries to decrypt it with the encryptionKey in the associated CryptoService.
// Returns the byte representation of the key or an error.
func (c *CryptoService) DecryptPrivateKey(key string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
