package utils

import "runtime"

const (
	baseChars     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	baseCharsSize = len(baseChars)

	Linux   = runtime.GOOS == "linux"
	Windows = runtime.GOOS == "windows"
)

var (
	baseBytesChars = []byte(baseChars)
)
