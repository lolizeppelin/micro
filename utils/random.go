package utils

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

var (
	RandomSeed *rand.Rand
)

func init() {
	RandomSeed = rand.New(rand.NewSource(time.Now().UnixNano())) // 初始化随机数种子
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
	return RandomSeed.Intn(n)
}

// RandString 返回一个指定长度的随机字符串
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = Base62Chars[rand.Intn(Base62CharsLen)]
	}
	return string(b)
}

// ShuffleSlice 使用 Fisher-Yates 算法随机打乱切片
func ShuffleSlice[T any](list []T) []T {
	shuffled := make([]T, len(list))
	copy(shuffled, list)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := RandomSeed.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
