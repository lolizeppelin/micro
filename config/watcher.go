package config

import (
	"context"
	"github.com/lolizeppelin/micro"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

func newMap() *wMap {
	return &wMap{
		stopped:  false,
		watchers: map[string]*watcher{},
	}
}

type wMap struct {
	sync.Mutex
	stopped  bool
	watchers map[string]*watcher
}

type watcher struct {
	key     string
	ch      clientv3.WatchChan
	exit    chan bool
	handler func(string, []*clientv3.Event, error)
}

func newWatcher(ec *EtcdConfig, key string, handler func(string, []*clientv3.Event, error)) {
	watchers := ec.watchers

	watchers.Lock()
	defer watchers.Unlock()
	if watchers.stopped {
		handler(key, nil, micro.ErrWatcherStopped)
		return
	}
	wc, ok := watchers.watchers[key]
	if ok {
		wc.Stop()
	}

	_w := &watcher{
		key:     key,
		ch:      ec.watcher.Watch(context.Background(), key, clientv3.WithPrefix()),
		exit:    make(chan bool),
		handler: handler,
	}
	watchers.watchers[key] = _w
	go _w.run()

}

func (w *watcher) run() {
	for {
		select {
		case rsp, ok := <-w.ch:
			if !ok {
				return
			}
			w.handler(w.key, rsp.Events, nil)
		case <-w.exit:
			w.handler(w.key, nil, micro.ErrWatcherStopped)
			return
		}
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		close(w.exit)
	}
	return nil
}

func stopAll(ec *EtcdConfig) {
	watchers := ec.watchers
	watchers.Lock()
	defer watchers.Unlock()
	watchers.stopped = true
	for _, wc := range watchers.watchers {
		wc.Stop()
	}
}
