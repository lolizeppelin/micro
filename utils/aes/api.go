package aes

import "errors"

func NewAES128(key string, iv string) (*CBC128PKCS7, error) {
	if len(key) != 16 || len(iv) != 16 {
		return nil, errors.New("key or iv length not 16")
	}
	_key := []byte(key)
	_iv := []byte(iv)
	_, err := newCBC128Decrypter(_key, _iv)
	if err != nil {
		return nil, err
	}
	_, err = newCBC128Encrypter(_key, _iv)
	if err != nil {
		return nil, err
	}
	return &CBC128PKCS7{key: _key, iv: _iv}, nil
}

type CBC128PKCS7 struct {
	key []byte
	iv  []byte
}

func (c *CBC128PKCS7) Encrypt(data []byte) []byte {
	crypto, _ := newCBC128Encrypter(c.key, c.iv)
	src := pkcs7Padding(data, crypto.BlockSize())
	dst := make([]byte, len(src))
	crypto.CryptBlocks(dst, src)
	return dst
}

func (c *CBC128PKCS7) Dencrypt(data []byte) ([]byte, error) {
	crypto, _ := newCBC128Decrypter(c.key, c.iv)
	//初始化解密数据接收切片
	dst := make([]byte, len(data))
	//执行解密
	crypto.CryptBlocks(dst, data)
	//去除填充
	result, err := pkcs7UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return result, nil

}
