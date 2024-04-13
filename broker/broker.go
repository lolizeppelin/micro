package broker

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(ctx context.Context, topic string, m *transport.Message) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error

// Message is a message send/received from the broker.

// Event is given to a subscription handler for processing.
type Event interface {
	Topic() string
	Message() *transport.Message
	Ack() error
	Error() error
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}
