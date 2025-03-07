//go:build !aix && !windows
// +build !aix,!windows

package systemd

import (
	"github.com/coreos/go-systemd/v22/daemon"
)

func SdNotify(state string) error {
	ok, err := daemon.SdNotify(false, state)
	if err != nil {
		return err
	}
	if !ok {
		return NotSupportSystemd
	}
	return nil
}
