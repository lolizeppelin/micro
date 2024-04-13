package selector

import (
	"github.com/lolizeppelin/micro"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random is a random strategy algorithm for node selection
func Random(services []*micro.Service) Next {
	nodes := make([]*micro.Node, 0, len(services))

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*micro.Node, error) {
		if len(nodes) == 0 {
			return nil, micro.ErrNoneServiceAvailable
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(services []*micro.Service) Next {
	nodes := make([]*micro.Node, 0, len(services))

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int()
	var mtx sync.Mutex

	return func() (*micro.Node, error) {
		if len(nodes) == 0 {
			return nil, micro.ErrNoneServiceAvailable
		}

		mtx.Lock()
		node := nodes[i%len(nodes)]
		i++
		mtx.Unlock()

		return node, nil
	}
}
