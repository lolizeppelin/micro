package config

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/tracing"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"sync"
)

type WatchHandler func(context.Context, string, []*clientv3.Event, error)

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
	ctx     context.Context
	key     string
	ch      clientv3.WatchChan
	exit    chan bool
	handler WatchHandler
}

func newWatcher(ctx context.Context, ec *EtcdConfig, key string, handler WatchHandler) {
	watchers := ec.watchers

	watchers.Lock()
	defer watchers.Unlock()
	if watchers.stopped {
		handler(ctx, key, nil, micro.ErrWatcherStopped)
		return
	}
	wc, ok := watchers.watchers[key]
	if ok {
		_ = wc.Stop()
	}
	_w := &watcher{
		ctx:     ctx,
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
			tracer := tracing.GetTracer(EtcdConfigScope, _version)
			ctx, span := tracer.Start(w.ctx, "etch.watcher",
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithAttributes(
					attribute.String("key", w.key),
				),
				oteltrace.WithAttributes(
					attribute.Int("count", len(rsp.Events)),
				),
			)
			defer span.End()

			w.handler(ctx, w.key, rsp.Events, nil)
		case <-w.exit:
			//ctx := w.ctx
			//w.handler(ctx, w.key, nil, micro.ErrWatcherStopped)
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
		_ = wc.Stop()
	}
}
