package registry

import (
	"fmt"
	"github.com/lolizeppelin/micro"
)

func addNodes(old, neu []*micro.Node) []*micro.Node {
	nodes := make([]*micro.Node, len(neu))
	// add all new nodes
	for i, n := range neu {
		node := *n
		nodes[i] = &node
	}

	// look at old nodes
	for _, o := range old {
		var exists bool

		// check against new nodes
		for _, n := range nodes {
			// ids match then skip
			if o.Id == n.Id {
				exists = true
				break
			}
		}

		// keep old node
		if !exists {
			node := *o
			nodes = append(nodes, &node)
		}
	}

	return nodes
}

func delNodes(old, del []*micro.Node) []*micro.Node {
	var nodes []*micro.Node
	for _, o := range old {
		var rem bool
		for _, n := range del {
			if o.Id == n.Id {
				rem = true
				break
			}
		}
		if !rem {
			nodes = append(nodes, o)
		}
	}
	return nodes
}

// CopyService make a copy of service
func CopyService(service *micro.Service) *micro.Service {
	// copy service
	s := new(micro.Service)
	*s = *service

	// copy nodes
	nodes := make([]*micro.Node, len(service.Nodes))
	for j, node := range service.Nodes {
		n := new(micro.Node)
		*n = *node
		nodes[j] = n
	}
	s.Nodes = nodes

	// copy endpoints
	eps := make([]*micro.Endpoint, len(service.Endpoints))
	for j, ep := range service.Endpoints {
		e := new(micro.Endpoint)
		*e = *ep
		eps[j] = e
	}
	s.Endpoints = eps
	return s
}

// CopyServices makes a copy of services
func CopyServices(current []*micro.Service) []*micro.Service {
	services := make([]*micro.Service, len(current))
	for i, service := range current {
		services[i] = CopyService(service)
	}
	return services
}

// Merge merges two lists of services and returns a new copy
func Merge(olist []*micro.Service, nlist []*micro.Service) []*micro.Service {
	var srv []*micro.Service

	for _, n := range nlist {
		var seen bool
		for _, o := range olist {
			if o.Version == n.Version {
				sp := new(micro.Service)
				// make copy
				*sp = *o
				// set nodes
				sp.Nodes = addNodes(o.Nodes, n.Nodes)

				// mark as seen
				seen = true
				srv = append(srv, sp)
				break
			} else {
				sp := new(micro.Service)
				// make copy
				*sp = *o
				srv = append(srv, sp)
			}
		}
		if !seen {
			srv = append(srv, CopyServices([]*micro.Service{n})...)
		}
	}
	return srv
}

// Remove removes services and returns a new copy
func Remove(old, del []*micro.Service) []*micro.Service {
	var services []*micro.Service

	for _, o := range old {
		srv := new(micro.Service)
		*srv = *o

		var rem bool

		for _, s := range del {
			if srv.Version == s.Version {
				srv.Nodes = delNodes(srv.Nodes, s.Nodes)

				if len(srv.Nodes) == 0 {
					rem = true
				}
			}
		}

		if !rem {
			services = append(services, srv)
		}
	}

	return services
}

// Topic 获取需要发送的主题
func Topic(service string, version *micro.Version, node string) string {
	var topic string
	if node != "" {
		topic = fmt.Sprintf("%s.node-%s.endpoints", service, node)
	} else {
		topic = fmt.Sprintf("%s.v%s.endpoints", service, version.Main())
	}
	return topic
}

func Topics(service string, version *micro.Version, node string) []string {
	t1 := Topic(service, version, "")
	t2 := Topic(service, version, node)
	return []string{t1, t2}
}
