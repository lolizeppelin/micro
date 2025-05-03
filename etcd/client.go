package etcd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net"
	"time"
)

func NewEtcdClient(opts ...Option) (*clientV3.Client, error) {

	options := Options{
		DialTimeout: 5 * time.Second,
	}

	cfg := clientV3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Logger:    zap.NewNop(),
	}
	for _, o := range opts {
		o(&options)
	}

	cfg.DialTimeout = options.DialTimeout
	if options.AuthCreds != nil {
		cfg.Username = options.AuthCreds.Username
		cfg.Password = options.AuthCreds.Password
	}

	if options.TLSConfig != nil {
		tlsConfig := options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		cfg.TLS = tlsConfig
	}

	var cAddress []string

	for _, address := range options.Address {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		var ae *net.AddrError
		if errors.As(err, &ae) && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddress = append(cAddress, net.JoinHostPort(addr, port))
		}
	}

	if len(cAddress) > 0 {
		cfg.Endpoints = cAddress
	}

	c, err := clientV3.New(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()
	_, err = c.Status(ctx, cfg.Endpoints[0])
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("connect etcd failed")
	}
	return c, err
}
