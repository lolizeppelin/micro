package utils

import "runtime"

const (
	Base62Chars    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Base62CharsLen = len(Base62Chars)

	Linux   = runtime.GOOS == "linux"
	Windows = runtime.GOOS == "windows"
)

var (
	baseBytesChars = []byte(Base62Chars)
)
