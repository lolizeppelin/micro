package registry

import (
	"context"
	"errors"
	"github.com/lolizeppelin/micro"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdWatcher struct {
	stop    chan bool
	w       clientv3.WatchChan
	client  *clientv3.Client
	timeout time.Duration
}

func newEtcdWatcher(r *etcdRegistry, timeout time.Duration, service string) (micro.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan bool, 1)

	go func() {
		<-stop
		cancel()
	}()

	watchPath := prefix
	if len(service) > 0 {
		watchPath = servicePath(service) + "/"
	}

	return &etcdWatcher{
		stop:    stop,
		w:       r.client.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV()),
		client:  r.client,
		timeout: timeout,
	}, nil
}

func (ew *etcdWatcher) Next() (*micro.Result, error) {
	for res := range ew.w {
		if res.Err() != nil {
			return nil, res.Err()
		}
		if res.Canceled {
			return nil, errors.New("could not get next")
		}
		for _, ev := range res.Events {
			service := decode(ev.Kv.Value)
			var action string

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"

				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}

			if service == nil {
				continue
			}
			return &micro.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}
	return nil, errors.New("could not get next")
}

func (ew *etcdWatcher) Stop() {
	select {
	case <-ew.stop:
		return
	default:
		close(ew.stop)
	}
}
