package micro

import (
	"errors"
)

var (
	// ErrIPNotFound no IP address found, and explicit IP not provided.
	ErrIPNotFound = errors.New("no IP address found, and explicit IP not provided")

	ErrServiceNotFound = errors.New("service not found")
	ErrWatcherStopped  = errors.New("watcher stopped")

	ErrSelectServiceNotFound  = errors.New("service not found")
	ErrSelectEndpointNotFound = errors.New("endpoint not found")
	ErrNoneServiceAvailable   = errors.New("none available") // node not found
)
