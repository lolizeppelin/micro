package utils

import (
	"bytes"
	"fmt"
	"github.com/sqids/sqids-go"
)

func ToBase62Bytes(n int) (s []byte) {
	if n <= 0 {
		s = []byte{'0'}
		return
	}
	for n > 0 {
		s = append(s, baseBytesChars[n%Base62CharsLen])
		n = n / Base62CharsLen
	}
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return
}

func ToBase62(n int) string {
	return string(ToBase62Bytes(n))
}

func FromBase62(s string) (n int, err error) {
	for _, c := range []byte(s) {
		index := bytes.IndexByte(baseBytesChars, c)
		if index < 0 {
			return -1, fmt.Errorf("string value error")
		}
		n = n*Base62CharsLen + index
	}
	return
}

func New62Sqid(len uint8) (*sqids.Sqids, error) {
	return sqids.New(sqids.Options{
		Alphabet:  Base62Chars,
		MinLength: len,
	})
}
