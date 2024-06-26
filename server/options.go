package server

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type options struct {
	//Broker broker.DefaultBroker,
	Id            uint32
	Name          string
	Address       string
	MaxMsgSize    int
	Version       *micro.Version // 当前服务版本号
	Interval      time.Duration
	Listener      net.Listener
	Broker        broker.Broker
	Registry      micro.Registry
	Components    []micro.Component
	GrpcOpts      []grpc.ServerOption
	BrokerOpts    []broker.SubscribeOption
	RegisterCheck func(context.Context) error
	WaitGroup     *sync.WaitGroup
	Metadata      map[string]string
}

type Option func(*options)

func WithServerId(id int32) Option {
	if id <= 0 || id > MaxServerSN {
		panic("server id over size")
	}
	return func(o *options) {
		o.Id = uint32(id)
	}
}

// WithName TODO 正则校验
func WithName(name string) Option {

	if name == "" {
		panic("name value error")
	}
	return func(o *options) {
		o.Name = name
	}
}

func WithVersion(version *micro.Version) Option {
	return func(o *options) {
		o.Version = version
	}
}

func WithAddress(address string) Option {
	return func(o *options) {
		o.Address = address
	}
}

func WithMaxMsgSize(size int) Option {
	if size <= 1024 {
		panic("grpc buff size error")
	}
	return func(o *options) {
		o.MaxMsgSize = size
	}
}

func WithRegisterCheckInterval(seconds time.Duration) Option {
	if seconds < 5 {
		panic("register check interval error")
	}
	return func(o *options) {
		o.Interval = seconds * time.Second
	}
}

func WithListener(listener net.Listener) Option {
	return func(o *options) {
		o.Listener = listener
	}
}

func WithBroker(broker broker.Broker) Option {
	return func(o *options) {
		o.Broker = broker
	}
}

func WithRegistry(registry micro.Registry) Option {
	return func(o *options) {
		o.Registry = registry
	}
}

func WithComponents(components ...micro.Component) Option {
	return func(o *options) {
		o.Components = components
	}
}

func WithGrpcOptions(opts ...grpc.ServerOption) Option {
	return func(o *options) {
		o.GrpcOpts = opts
	}
}

func WithRegisterCheck(f func(context.Context) error) Option {
	return func(o *options) {
		o.RegisterCheck = f
	}
}

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(o *options) {
		o.WaitGroup = wg
	}
}

// WithMetadata associated with the server
func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.Metadata = md
	}
}
