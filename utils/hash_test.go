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

	fmt.Printf("hash 1 to %d\n", HashStringToInt("0", 0, 3))
	fmt.Printf("hash 1 to %d\n", HashStringToInt("1", 0, 3))
	fmt.Printf("hash 2 to %d\n", HashStringToInt("2", 0, 3))
	fmt.Printf("hash 3 to %d\n", HashStringToInt("3", 0, 3))
	fmt.Printf("hash 4 to %d\n", HashStringToInt("4", 0, 3))
	fmt.Printf("hash 5 to %d\n", HashStringToInt("5", 0, 3))
	fmt.Printf("hash 6 to %d\n", HashStringToInt("6", 0, 3))
	fmt.Printf("hash 6 to %d\n", HashStringToInt("7", 0, 3))
	fmt.Printf("hash 6 to %d\n", HashStringToInt("88", 0, 3))

}
