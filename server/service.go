package grpc

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
)

type Service struct {
	opts      *options
	services  map[string]map[string]*CompHandler
	endpoints []*micro.Endpoint
	registry  *micro.Service
	//subscribers map[string]broker.Handler
	subscribed map[string]broker.Subscriber
}

func (s *Service) Handler(service string, method string) *CompHandler {
	sv, ok := s.services[service]
	if !ok {
		return nil
	}
	return sv[method]
}

func newService(opts *options) *Service {

	services := ExtractComponents(opts.Components)
	endpoints := extractEndpoints(services)

	node := &micro.Node{
		Id:       SNBase62(opts.Id),
		Version:  opts.Version.Version(), // 节点版本号
		Address:  opts.Address,
		Metadata: opts.Metadata,
	}

	reg := &micro.Service{
		Name:      opts.Name,
		Version:   opts.Version.Main(), // 服务主版本号
		Nodes:     []*micro.Node{node},
		Endpoints: endpoints,
		//Metadata:
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
		registry:  reg,
	}

}
