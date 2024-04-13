package registry

import (
	"crypto/tls"
	"time"
)

type Options struct {
	TLSConfig *tls.Config
	Address   []string
	Timeout   time.Duration
	Secure    bool
	TTL       time.Duration
}

type Option func(*Options)
