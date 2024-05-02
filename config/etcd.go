package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
)

type EtcdConfig struct {
	prefix   string
	client   *clientv3.Client
	kv       clientv3.KV
	watcher  clientv3.Watcher
	watchers *wMap
}

func NewEtcdConfig(opts ...Option) (*EtcdConfig, error) {

	options := NewOptions(opts...)

	var endpoints []string

	for _, a := range options.Address {
		addr, port, err := net.SplitHostPort(a)
		var ae *net.AddrError
		if errors.As(err, &ae) && ae.Err == "missing port in address" {
			port = "2379"
			addr = a
			endpoints = append(endpoints, fmt.Sprintf("%s:%s", addr, port))
		}
	}

	if len(endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}

	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: options.Timeout,
	}

	if options.AuthCreds != nil {
		cfg.Username = options.AuthCreds.Username
		cfg.Password = options.AuthCreds.Password
	}

	// use default config
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	c := clientv3.NewKV(client)
	w := clientv3.NewWatcher(client)

	return &EtcdConfig{
		prefix:  options.Prefix,
		client:  client,
		kv:      c,
		watcher: w,
	}, nil

}

func (e *EtcdConfig) Get(ctx context.Context, key string) (*mvccpb.KeyValue, error) {
	kv, err := e.kv.Get(ctx, e.prefix+key)
	if err != nil {
		return nil, err
	}
	if kv.Count < 1 {
		return nil, micro.ErrConfigFound
	}
	if kv.Count > 1 {
		return nil, micro.ErrConfigCountNotMatch
	}
	return kv.Kvs[0], nil
}

func (e *EtcdConfig) List(ctx context.Context, key string) ([]*mvccpb.KeyValue, error) {
	kv, err := e.kv.Get(ctx, e.prefix+key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if kv.Count < 1 {
		return nil, micro.ErrConfigFound
	}
	return kv.Kvs, nil
}

func (e *EtcdConfig) Put(ctx context.Context, key string, value string) (*mvccpb.KeyValue, error) {
	kv, err := e.kv.Put(ctx, e.prefix+key, value)
	if err != nil {
		return nil, err
	}
	return kv.PrevKv, nil
}

func (e *EtcdConfig) Delete(ctx context.Context, key string) (int64, error) {
	kv, err := e.kv.Delete(ctx, e.prefix+key)
	if err != nil {
		return 0, err
	}
	return kv.Deleted, nil
}

func (e *EtcdConfig) Truncate(ctx context.Context, key string) (int64, error) {
	kv, err := e.kv.Delete(ctx, e.prefix+key, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return kv.Deleted, nil
}

func (e *EtcdConfig) Watch(key string, handler func([]*clientv3.Event, error)) {
	newWatcher(e, e.prefix+key, handler)
}

func (e *EtcdConfig) Close() error {
	stopAll(e)
	return e.client.Close()
}
