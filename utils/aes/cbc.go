package aes

import (
	"crypto/aes"
	"crypto/cipher"
)

func newCBC128Decrypter(key []byte, iv []byte) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ins := cipher.NewCBCDecrypter(block, iv)
	return ins, nil
}

func newCBC128Encrypter(key []byte, iv []byte) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ins := cipher.NewCBCEncrypter(block, iv)
	return ins, nil
}
