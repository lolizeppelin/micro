// Package secrets is an interface for encrypting and decrypting secrets
package secrets

// Secrets encrypts or decrypts arbitrary data. The data should be as small as possible.
type Secrets interface {
	// Initialize options
	Init(...Option) error
	// Return the options
	Options() Options
	// Decrypt a value
	Decrypt([]byte, ...DecryptOption) ([]byte, error)
	// Encrypt a value
	Encrypt([]byte, ...EncryptOption) ([]byte, error)
	// Secrets implementation
	String() string
}

// DecryptOptions can be passed to Secrets.Decrypt.
