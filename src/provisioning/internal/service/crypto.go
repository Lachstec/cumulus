package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
)

type CryptoService struct {
	encryptionKey []byte
}

func NewCryptoService(encryptionKey []byte) *CryptoService {
	return &CryptoService{encryptionKey: encryptionKey}
}

func (c *CryptoService) encryptPrivateKey(key ed25519.PrivateKey) (string, error) {
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
