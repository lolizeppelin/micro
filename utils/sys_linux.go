//go:build !aix && !windows
// +build !aix,!windows

package utils

import (
	"os"
	"syscall"
)

func isUnixSocketFile(stat os.FileInfo) bool {
	st, ok := stat.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	socketMode := st.Mode & syscall.S_IFMT
	return socketMode == syscall.S_IFSOCK
}
