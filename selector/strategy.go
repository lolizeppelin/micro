package selector

import (
	"github.com/lolizeppelin/micro"
	"github.com/minio/highwayhash"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// zeroKey is the base key for all hashes, it is 32 zeros.
var zeroKey [32]byte

// NewSharedStrategy returns a `SelectOption` that directs all request according to the given `keys`.
func NewSharedStrategy(keys []string) Strategy {
	return func(services []*micro.Service) Next {
		return next(keys, services)
	}
}

// Next returns a `Next` function which returns the next highest scoring node.
func next(keys []string, services []*micro.Service) Next {
	possibleNodes, scores := ScoreNodes(keys, services)

	return func() (*micro.Node, error) {
		var best uint64
		pos := -1

		// Find the best scoring node from those available.
		for i, score := range scores {
			if score >= best && possibleNodes[i] != nil {
				best = score
				pos = i
			}
		}

		if pos < 0 {
			// There was no node found.
			return nil, micro.ErrNoneServiceAvailable
		}

		// Choose this node and set it's score to zero to stop it being selected again.
		node := possibleNodes[pos]
		possibleNodes[pos] = nil
		scores[pos] = 0
		return node, nil
	}
}

// ScoreNodes returns a score for each node found in the given services.
func ScoreNodes(keys []string, services []*micro.Service) (possibleNodes []*micro.Node, scores []uint64) {
	// Generate a base hashing key based off the supplied keys values.
	key := highwayhash.Sum([]byte(strings.Join(keys, ":")), zeroKey[:])

	// Get all the possible nodes for the services, and assign a hash-based score to each of them.
	for _, s := range services {
		for _, n := range s.Nodes {
			// Use the base key from above to calculate a derivative 64 bit hash number based off the instance ID.
			score := highwayhash.Sum64([]byte(n.Id), key[:])
			scores = append(scores, score)
			possibleNodes = append(possibleNodes, n)
		}
	}
	return
}

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
