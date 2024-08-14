package systemd

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/daemon"
)

// ErrSdNotifyNoSocket is the error returned when the NOTIFY_SOCKET does not exist.

// Ready sends READY=1 to the systemd notify socket.
func Ready() error {
	return SdNotify(daemon.SdNotifyReady)
}

// Stopping sends STOPPING=1 to the systemd notify socket.
func Stopping() error {
	return SdNotify(daemon.SdNotifyStopping)
}

// Reloading sends RELOADING=1 to the systemd notify socket.
func Reloading() error {
	return SdNotify(daemon.SdNotifyReloading)
}

// Errno sends ERRNO=? to the systemd notify socket.
func Errno(errno int) error {
	return SdNotify(fmt.Sprintf("ERRNO=%d", errno))
}

// Status sends STATUS=? to the systemd notify socket.
func Status(status string) error {
	return SdNotify("STATUS=" + status)
}

// Watchdog sends WATCHDOG=1 to the systemd notify socket.
func Watchdog() error {
	return SdNotify(daemon.SdNotifyWatchdog)
}
