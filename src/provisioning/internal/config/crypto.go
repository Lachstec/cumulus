package config

// CryptoConfig contains configuration Options for everything related
// to cryptography in the application
type CryptoConfig struct {
	// EncryptionKey that is used to securely store SSH keys for the gameservers in the DB
	EncryptionKey []byte
}
