package utils

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	a := HashString(MD5, "1111", "222")
	b := HashString(SHA1, "1111", "222")
	c := HashString(SHA256, "1111", "222")
	d := HashString(MD5, "1111", "222")
	e := HashString(CRC32, "12312312311231231", "222")

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)

	fmt.Println(Md5Hmac("afaga", "agaga"))
	fmt.Println(Sha1Hmac("afaga", "agaga"))
	fmt.Println(Sha256Hmac("afaga", "agaga"))
}
