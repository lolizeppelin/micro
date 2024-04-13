package utils

import (
	"os"
)

func isUnixSocketFile(stat os.FileInfo) bool {
	return false
}
