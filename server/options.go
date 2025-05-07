package server

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
	"github.com/lolizeppelin/micro/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sync"
	"time"
)

type Options struct {
	Id            uint64
	Name          string
	MaxMsgSize    int
	Version       *micro.Version // 当前服务版本号
	Min           *micro.Version // 支持的最小版本(默认当前版本)
	Max           *micro.Version // 支持的最大版本(默认当前版本)
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

	Credentials credentials.TransportCredentials
}

type Option func(*Options)

func WithServerId(id uint64) Option {
	if id == 0 || id > MaxServerSN {
		panic("server id over size")
	}
	return func(o *Options) {
		o.Id = id
	}
}

// WithName TODO 正则校验
func WithName(name string) Option {
	if name == "" {
		panic("name value error")
	}
	return func(o *Options) {
		o.Name = name
	}
}

func WithVersion(version *micro.Version) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func WithMax(version *micro.Version) Option {
	return func(o *Options) {
		if o.Version != nil && o.Version.Compare(*version) >= 0 {
			panic("max version value error")
		}
		o.Max = version
	}
}

func WithMin(version *micro.Version) Option {
	return func(o *Options) {
		if o.Version != nil && o.Version.Compare(*version) <= 0 {
			panic("min version value error")
		}
		o.Min = version
	}
}

func WithMaxMsgSize(size int) Option {
	if size <= 1024 {
		panic("grpc buff size error")
	}
	return func(o *Options) {
		o.MaxMsgSize = size
	}
}

func WithRegisterCheckInterval(seconds time.Duration) Option {
	if seconds < 5 {
		panic("register check interval error")
	}
	return func(o *Options) {
		o.Interval = seconds * time.Second
	}
}

func WithListener(listener net.Listener) Option {
	return func(o *Options) {
		o.Listener = listener
	}
}

func WithBroker(broker broker.Broker) Option {
	return func(o *Options) {
		o.Broker = broker
	}
}

func WithRegistry(registry micro.Registry) Option {
	return func(o *Options) {
		o.Registry = registry
	}
}

func WithComponents(components ...micro.Component) Option {
	return func(o *Options) {
		o.Components = components
	}
}

func WithGrpcOptions(opts ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.GrpcOpts = opts
	}
}

func WithRegisterCheck(f func(context.Context) error) Option {
	return func(o *Options) {
		o.RegisterCheck = f
	}
}

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(o *Options) {
		o.WaitGroup = wg
	}
}

// WithMetadata associated with the server
func WithMetadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// WithCredentials 设置证书
func WithCredentials(credentials credentials.TransportCredentials) Option {
	return func(o *Options) {
		o.Credentials = credentials
	}
}

func WithBrokerOpts(options []broker.SubscribeOption) Option {
	return func(o *Options) {
		o.BrokerOpts = options
	}
}

func NewOptions(name string) *Options {
	return &Options{
		Name:          name,
		MaxMsgSize:    DefaultMaxMsgSize,
		Interval:      time.Second * 30,
		RegisterCheck: registry.DefaultRegisterCheck,
		Credentials:   insecure.NewCredentials(),
		WaitGroup:     new(sync.WaitGroup),
	}

}
