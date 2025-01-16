package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/ssh"
)

// CryptoService provides utility functions to handle SSH keys
// that are used in order to execute commands on the game servers.
// This service should not be directly exposed to end user requests.
type CryptoService struct {
	// encryptionKey is used to encrypt SSH keys in order to serverstore them in the database securely.
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

// publicKeyToOpenSSH is a helper function that turns an ed25519 Public Key into
// a string to give to OpenStack.
func publicKeyToOpenSSH(pub ed25519.PublicKey) (string, error) {
	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return "", err
	}

	publicKeyString := "ssh-ed25519" + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())
	return publicKeyString, nil
}

// NewKeyPair generates a new Key Pair for a gameserver. It returns
// the private key as an encrypted, base64-encoded string and the public key as
// a string compatible with OpenSSH / Openstack. Returns an error if something goes wrong.
func (c *CryptoService) NewKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	publicKey, err := publicKeyToOpenSSH(pub)
	if err != nil {
		return "", "", err
	}

	privateKey, err := c.EncryptPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	return publicKey, privateKey, nil
}
