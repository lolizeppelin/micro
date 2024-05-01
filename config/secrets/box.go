// Package box is an asymmetric implementation of config/secrets using nacl/box
package secrets

import (
	"crypto/rand"

	"github.com/pkg/errors"
	naclbox "golang.org/x/crypto/nacl/box"
)

type box struct {
	options Options

	publicKey  [keyLength]byte
	privateKey [keyLength]byte
}

// NewSecrets returns a nacl-box codec.
func NewSecrets(opts ...Option) Secrets {
	b := &box{}
	for _, o := range opts {
		o(&b.options)
	}
	return b
}

func (b *box) Init(opts ...Option) error {
	for _, o := range opts {
		o(&b.options)
	}
	if len(b.options.PrivateKey) != keyLength || len(b.options.PublicKey) != keyLength {
		return errors.Errorf("a public key and a private key of length %d must both be provided", keyLength)
	}
	copy(b.privateKey[:], b.options.PrivateKey)
	copy(b.publicKey[:], b.options.PublicKey)
	return nil
}

// Options returns options.
func (b *box) Options() Options {
	return b.options
}

// String returns nacl-box.
func (*box) String() string {
	return "nacl-box"
}

// Encrypt encrypts a message with the sender's private key and the receipient's public key.
func (b *box) Encrypt(in []byte, opts ...EncryptOption) ([]byte, error) {
	var options EncryptOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.RecipientPublicKey) != keyLength {
		return []byte{}, errors.New("recepient's public key must be provided")
	}
	var recipientPublicKey [keyLength]byte
	copy(recipientPublicKey[:], options.RecipientPublicKey)
	var nonce [24]byte
	if _, err := rand.Reader.Read(nonce[:]); err != nil {
		return []byte{}, errors.Wrap(err, "couldn't obtain a random nonce from crypto/rand")
	}
	return naclbox.Seal(nonce[:], in, &nonce, &recipientPublicKey, &b.privateKey), nil
}

// Decrypt Decrypts a message with the receiver's private key and the sender's public key.
func (b *box) Decrypt(in []byte, opts ...DecryptOption) ([]byte, error) {
	var options DecryptOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.SenderPublicKey) != keyLength {
		return []byte{}, errors.New("sender's public key bust be provided")
	}
	var nonce [24]byte
	var senderPublicKey [32]byte
	copy(nonce[:], in[:24])
	copy(senderPublicKey[:], options.SenderPublicKey)
	decrypted, ok := naclbox.Open(nil, in[24:], &nonce, &senderPublicKey, &b.privateKey)
	if !ok {
		return []byte{}, errors.New("incoming message couldn't be verified / decrypted")
	}
	return decrypted, nil
}
