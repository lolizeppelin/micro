package server

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/log"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
	"github.com/lolizeppelin/micro/utils"
	"net"
	"sync"
	"time"

	"github.com/lolizeppelin/micro/registry"

	"google.golang.org/grpc"
)

var (
	// DefaultMaxMsgSize define maximum message size that server can send
	// or receive.  Default value is 4MB.
	DefaultMaxMsgSize = 1024 * 1024 * 4
)

type RPCServer struct {
	tp.UnimplementedTransportServer

	sync.RWMutex
	wg   *sync.WaitGroup
	exit chan chan error

	started    bool
	registered bool

	opts    *options
	service *Service

	// grpc server
	server *grpc.Server
}

func newGRPCServer(opts *options) *RPCServer {

	// create a grpc server
	srv := &RPCServer{
		opts:    opts,
		exit:    make(chan chan error),
		service: newService(opts),
		wg:      opts.WaitGroup,
	}
	// configure the grpc server

	_opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(opts.MaxMsgSize),
		grpc.MaxSendMsgSize(opts.MaxMsgSize),
		//grpc.UnknownServiceHandler(srv.handler),
		grpc.ConnectionTimeout(10 * time.Second),
	}

	_opts = append(_opts, opts.GrpcOpts...)
	srv.server = grpc.NewServer(_opts...)

	tp.RegisterTransportServer(srv.server, srv)

	return srv
}

func (g *RPCServer) Start() error {
	g.RLock()
	if g.started {
		g.RUnlock()
		return nil
	}
	g.RUnlock()

	config := g.opts

	log.Infof("Server [grpc] Listening on %s", g.opts.Address)

	go func() {
		if err := g.server.Serve(g.opts.Listener); err != nil {
			log.Errorf("gRPC Server start error: %v", err)
		}
	}()

	go func() {
		if err := g.Register(); err != nil {
			log.Errorf("Server register error: %s", err.Error())
		}
	}()

	go func() {
		t := time.NewTicker(g.opts.Interval)

		// return error chan
		var (
			err error
			ch  chan error
		)

	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				g.RLock()
				registered := g.registered
				g.RUnlock()
				ctx := context.Background()

				checkErr := g.opts.RegisterCheck(ctx)
				if checkErr != nil && registered {
					log.Errorf("Server %s-%d register check error: %s, deregister it",
						config.Name, config.Id, checkErr.Error())
					// deregister self in case of error
					if err = g.Deregister(); err != nil {
						log.Errorf("Server %s-%s deregister error: %s", config.Name, config.Id, err)
					}
				} else if checkErr != nil && !registered {
					log.Errorf("Server %s-%d register check error: %s",
						config.Name, config.Id, checkErr.Error())
					continue
				}
				// Register 内部包含续租
				if err = g.Register(); err != nil {
					log.Errorf("Server register error: %s", err.Error())
				}

			// wait for exit
			case ch = <-g.exit:
				break Loop
			}
		}

		// deregister self
		if err = g.Deregister(); err != nil {
			log.Errorf("server deregister error: %s", err.Error())
		}
		// wait for waitgroup
		g.wg.Wait()
		// stop the grpc server
		exit := make(chan bool)

		go func() {
			g.server.GracefulStop()
			close(exit)
		}()

		select {
		case <-exit:
		case <-time.After(time.Second):
			g.server.Stop()
		}

		if config.Broker != nil {
			log.Infof("broker [%s] Disconnected from %s", config.Broker.String(), config.Broker.Address())
			// disconnect broker
			if err = config.Broker.Disconnect(); err != nil {
				log.Errorf("broker [%s] disconnect error: %v", config.Broker.String(), err)
			}
		}

		// close transport
		ch <- err
	}()

	// mark the server as started
	g.Lock()
	g.started = true
	g.Unlock()

	return nil
}

func (g *RPCServer) Stop() error {

	g.RLock()
	if !g.started {
		g.RUnlock()
		return nil
	}
	g.RUnlock()

	ch := make(chan error)
	g.exit <- ch

	var err error
	select {
	case err = <-ch:
		g.Lock()
		g.server = nil
		g.started = false
		g.Unlock()
	}
	return err
}

func NewServer(o ...Option) (*RPCServer, error) {

	opts := &options{
		Id:            1,
		Name:          "server",
		MaxMsgSize:    DefaultMaxMsgSize,
		Interval:      time.Second * 30,
		RegisterCheck: registry.DefaultRegisterCheck,
	}

	for _, f := range o {
		f(opts)
	}

	if opts.Listener == nil {
		ls, err := net.Listen("tcp", "127.0.0.1:1780")
		if err != nil {
			return nil, err
		}
		opts.Listener = ls
	}

	if opts.Address == "" {
		address := opts.Listener.Addr().String()
		opts.Address = address
	}

	if !utils.VerifyAddr(opts.Address) {
		//return nil, fmt.Errorf("address value error")
	}

	if opts.WaitGroup == nil {
		opts.WaitGroup = new(sync.WaitGroup)
	}

	if opts.Version == nil {
		version, _ := micro.NewVersion("1.0")
		opts.Version = version
	}

	if opts.Metadata == nil {
		opts.Metadata = map[string]string{}
	}

	return newGRPCServer(opts), nil
}
