package server

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
	"github.com/lolizeppelin/micro/utils"
)

type Service struct {
	opts       *Options
	services   map[string]map[string]*Handler
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

	node.Metadata["registry"] = opts.Registry.Name()
	if opts.Broker != nil {
		node.Metadata["broker"] = opts.Broker.Name()
	}
	//node.Metadata["server"] = g.String()
	//node.Metadata["transport"] = g.String()
	node.Metadata["protocol"] = "grpc"

	emap, err := utils.SliceToMapByField[*micro.Endpoint, string](endpoints, "Name")
	if err != nil {
		panic("convert endpoint list to map failed")
	}

	return &Service{
		opts:     opts,
		services: services,
		//endpoints:  endpoints,
		subscribed: map[string]broker.Subscriber{},
		registry: &micro.Service{
			Name:      opts.Name,
			Version:   opts.Version.Major,
			Nodes:     []*micro.Node{node},
			Endpoints: emap,
			Metadata:  map[string]string{},
		},
	}

}
