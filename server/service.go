package server

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
)

type Service struct {
	opts       *Options
	services   map[string]map[string]*Handler
	endpoints  []*micro.Endpoint
	registry   *micro.Service
	subscribed map[string]broker.Subscriber
}

func (s *Service) Handler(service string, method string) *Handler {
	sv, ok := s.services[service]
	if !ok {
		return nil
	}
	return sv[method]
}

func newService(opts *Options) *Service {

	services, _ := ExtractComponents(opts.Components)
	endpoints := extractEndpoints(services)

	node := &micro.Node{
		Id:       SNBase62(opts.Id),
		Version:  *opts.Version, // 节点版本号
		Max:      opts.Max,
		Min:      opts.Min,
		Address:  opts.Listener.Addr().String(),
		Metadata: opts.Metadata,
	}

	//node.Metadata["broker"] = config.Broker.String()
	//node.Metadata["registry"] = config.Registry.String()
	//node.Metadata["server"] = g.String()
	//node.Metadata["transport"] = g.String()
	node.Metadata["protocol"] = "grpc"

	return &Service{
		opts:      opts,
		services:  services,
		endpoints: endpoints,
		registry: &micro.Service{
			Name:      opts.Name,
			Version:   opts.Version.Major,
			Nodes:     []*micro.Node{node},
			Endpoints: endpoints,
			Metadata:  map[string]string{},
		},
	}

}
