package utils

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())) // 初始化随机数种子
}

func RandomHex(length int) string {
	l := length + 2
	b := make([]byte, l)
	_, err := cr.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)[2:l]
}

func RandomInt(n int) int {
	return rand.Intn(n)
}

// RandString 返回一个指定长度的随机字符串
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = baseChars[rand.Intn(baseCharsSize)]
	}
	return string(b)
}
