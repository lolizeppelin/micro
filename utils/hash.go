package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"strings"
)

type HashMode string

const (
	MD5    HashMode = "md5"
	SHA1   HashMode = "sha1"
	SHA256 HashMode = "sha256"
	CRC32  HashMode = "crc32"
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

func HashString(mode HashMode, s ...string) string {
	buff := []byte(strings.Join(s, ""))
	switch mode {
	case "md5":
		return MD5Sum(buff)
	case "sha1":
		return Sha1Sum(buff)
	case "sha256":
		return Sha256Sum(buff)
	case "crc32":
		return CRC32Sum(buff)
	default:
		panic("error hash type")
	}
}

func Sha256Sum(value []byte) string {
	b := sha256.Sum256(value)
	return hex.EncodeToString(b[:])
}

func MD5Sum(value []byte) string {
	hash := md5.Sum(value)
	return hex.EncodeToString(hash[:])
}

func Sha1Sum(value []byte) string {
	h := sha1.New()
	h.Write(value)
	return hex.EncodeToString(h.Sum(nil))
}

func CRC32Sum(value []byte) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE(value))
}

func Sha1Hmac(value, salt string) string {
	h := hmac.New(sha1.New, []byte(salt))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Sha256Hmac(value, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Md5Hmac(value, salt string) string {
	h := hmac.New(md5.New, []byte(salt))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
