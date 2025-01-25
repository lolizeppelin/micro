package config

import (
	"context"
	"github.com/lolizeppelin/micro"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

const (
	DefaultPrefix = "/micro/config/"
)

type EtcdConfig struct {
	prefix   string
	kv       clientv3.KV
	watcher  clientv3.Watcher
	watchers *wMap
}

func NewEtcdConfig(client *clientv3.Client, prefix ...string) *EtcdConfig {
	p := DefaultPrefix
	if len(prefix) > 0 && prefix[0] != "" {
		p = prefix[0]
	}
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}

	c := clientv3.NewKV(client)
	w := clientv3.NewWatcher(client)

	return &EtcdConfig{
		prefix: p,
		//client:  client,
		kv:       c,
		watcher:  w,
		watchers: newMap(),
	}

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

/*
	SafePut

版本号小于等于指定版本号
*/
func (e *EtcdConfig) SafePut(ctx context.Context, key string, value string, version int64) (int64, error) {
	_key := e.prefix + key
	resp, err := e.kv.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(_key), "<=", version)).
		Then(clientv3.OpPut(_key, value)).
		Commit()
	if err != nil {
		return 0, err
	}
	if !resp.Succeeded {
		return 0, micro.ErrResultFailed
	}

	for _, res := range resp.Responses {
		result := res.GetResponsePut()
		if result != nil {
			return result.Header.Revision, nil
		}
	}

	return 0, micro.ErrUnknown
}

/*
	UpdatePut

版本号等于指定版本号
*/
func (e *EtcdConfig) UpdatePut(ctx context.Context, key string, value string, version int64) (int64, error) {

	_key := e.prefix + key
	resp, err := e.kv.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(_key), "=", version)).
		Then(clientv3.OpPut(_key, value)).
		Commit()
	if err != nil {
		return 0, err
	}
	if !resp.Succeeded {
		return 0, micro.ErrResultFailed
	}

	for _, res := range resp.Responses {
		result := res.GetResponsePut()
		if result != nil {
			return result.Header.Revision, nil
		}
	}

	return 0, micro.ErrUnknown
}

func (e *EtcdConfig) Watch(key string, handler func(string, []*clientv3.Event, error)) {
	newWatcher(e, e.prefix+key, handler)
}

func (e *EtcdConfig) Close() {
	stopAll(e)
}
