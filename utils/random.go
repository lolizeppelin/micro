package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
}

func RandomHex(length int) string {
	l := length + 2
	b := make([]byte, l)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2:l]
}

func RandomInt(n int) int {
	return rand.Intn(n)
}
