package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func DJB33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}

func Sha256Hash(value string, salt string) string {
	b := sha256.Sum256([]byte(value))
	b = sha256.Sum256([]byte(hex.EncodeToString(b[:]) + salt))
	return hex.EncodeToString(b[:])
}

func MD5Sum(value string, salt string) string {
	hash := md5.Sum([]byte(value + salt))
	return hex.EncodeToString(hash[:])
}

func Sha1Sum(value string, salt string) string {
	h := sha1.New()
	h.Write([]byte(value))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1Hmac(value string, salt string) string {
	h := hmac.New(sha1.New, []byte(salt))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Md5Hmac(value string, salt string) string {
	h := hmac.New(md5.New, []byte(salt))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
