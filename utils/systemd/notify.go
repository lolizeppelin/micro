//go:build !linux
// +build !linux

package systemd

// SdNotify sends a specified string to the systemd notification socket.
func SdNotify(state string) error {
	// do nothing
	return nil
}
