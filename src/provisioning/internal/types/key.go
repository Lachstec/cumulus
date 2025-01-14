package types

// Key represents a SSH Keypair for use with OpenStack.
type Key struct {
	// Name of the Keypair in OpenStack
	Name string
	// PublicKey part of the Keypair.
	PublicKey []byte
	// PrivateKey part of the Keypair.
	PrivateKey []byte
}
