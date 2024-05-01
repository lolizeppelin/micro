// Package source is the interface for sources
package source

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrWatcherStopped is returned when source watcher has been stopped.
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Source is the source from which config is loaded.
type Source interface {
	Read() (*ChangeSet, error)
	Write(*ChangeSet) error
	Watch() (Watcher, error)
	String() string
}

// ChangeSet represents a set of changes from a source.
type ChangeSet struct {
	Timestamp time.Time
	Checksum  string
	Format    string
	Source    string
	Data      []byte
}

func (c *ChangeSet) Sum() string {
	h := md5.New()
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Watcher watches a source for changes.
type Watcher interface {
	Next() (*ChangeSet, error)
	Stop() error
}
