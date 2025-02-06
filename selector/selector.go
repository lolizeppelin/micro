package selector

import "github.com/lolizeppelin/micro"

type Selector interface {
	// Select returns a function which should return the next node
	Select(service string, filters ...Filter) (Next, error)
	// Mark sets the success/error against a node
	Mark(service string, node *micro.Node, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error

	Name() string
}

// Next is a function that returns the next node
// based on the selector's strategy.
type Next func() (*micro.Node, error)

// Filter is used to filter a service during the selection process.
type Filter func([]*micro.Service) ([]*micro.Service, error)

// Strategy is a selection strategy e.g random, round robin.
type Strategy func([]*micro.Service) Next
