package services

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func generateTestEncryptionKey() []byte {
	encryptionKey := make([]byte, 32)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		panic(fmt.Sprintf("failed to generate encryption key: %v", err))
	}
	return encryptionKey
}

func TestEncryptDecryptPrivateKey(t *testing.T) {
	encryptionKey := generateTestEncryptionKey()
	cryptoService := NewCryptoService(encryptionKey)

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	encryptedPrivateKey, err := cryptoService.EncryptPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to encrypt private key: %v", err)
	}

	decryptedPrivateKey, err := cryptoService.DecryptPrivateKey(encryptedPrivateKey)
	if err != nil {
		t.Fatalf("Failed to decrypt private key: %v", err)
	}

	assert.Equal(t, privateKey.Seed(), decryptedPrivateKey, "Decrypted private key should match the original private key")
}

func TestPublicKeyToOpenSSH(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	sshKey, err := publicKeyToOpenSSH(publicKey)
	if err != nil {
		t.Fatalf("Failed to convert public key to OpenSSH: %v", err)
	}

	if !strings.HasPrefix(sshKey, "ssh-ed25519") {
		t.Fatalf("Public key does not start with 'ssh-ed25519': %s", sshKey)
	}
}

func TestNewKeyPair(t *testing.T) {
	encryptionKey := generateTestEncryptionKey()
	cryptoService := NewCryptoService(encryptionKey)

	publicKey, privateKey, err := cryptoService.NewKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate new key pair: %v", err)
	}

	if !strings.HasPrefix(publicKey, "ssh-ed25519") {
		t.Fatalf("Public key does not start with 'ssh-ed25519': %s", publicKey)
	}

	decryptedPrivateKey, err := cryptoService.DecryptPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to decrypt private key: %v", err)
	}

	assert.NotNil(t, decryptedPrivateKey, "Decrypted private key should not be nil")
}
