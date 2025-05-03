package broker

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/transport"
)

const (
	HandlerScope = "micro/broker/handler"
)

var (
	_version, _ = micro.NewVersion("1.0.0")
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Connect() error
	Disconnect() error
	Publish(ctx context.Context, topic string, m *transport.Message) error
	Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	Name() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(ctx context.Context, event Event) error

// Message is a message send/received from the broker.

// Event is given to a subscription handler for processing.
type Event interface {
	Message() *transport.Message
	Ack() error
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Topic() string
	Unsubscribe() error
}

/* ------------- event -------------*/

type kafkaEvent struct {
	msg *transport.Message
}

func (s *kafkaEvent) Message() *transport.Message {
	return s.msg
}

func (s *kafkaEvent) Ack() error {
	return nil
}
