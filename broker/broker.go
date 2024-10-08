package broker

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/vmihailenco/msgpack/v5"
)

type KafkaBroker struct {
	producer *kgo.Client
	opts     *Options
}

func (k *KafkaBroker) Connect() error {
	if k.producer != nil {
		return fmt.Errorf("producer connected")
	}
	producer, err := NewKafkaProducer(k.opts.Address)
	if err != nil {
		return err
	}
	k.producer = producer
	return nil
}

func (k *KafkaBroker) Disconnect() error {
	k.producer.Close()
	return nil
}

func (k *KafkaBroker) Publish(ctx context.Context, topic string, msg *transport.Message) error {
	if k.producer == nil {
		return fmt.Errorf("producer not connect")
	}
	buff, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}
	record := &kgo.Record{
		Topic: topic,
		Value: buff,
	}
	if key := ctx.Value(_ctxKey); key != nil {
		record.Key = []byte(key.(string))
	}
	if key := ctx.Value(_ctxPartKey); key != nil {
		record.Partition = key.(int32)
	}

	k.producer.TryProduce(ctx, record, func(record *kgo.Record, err error) {
		k.opts.ErrorHandler(0, record, err)
	})
	return nil
}

func (k *KafkaBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions(opts...)
	client, err := NewKafkaConsumer(k.opts.Address, topic, options)
	if err != nil {
		return nil, err
	}
	subscriber := &KafkaSubscriber{
		topic:     topic,
		client:    client,
		handler:   handler,
		fallback:  k.opts.ErrorHandler,
		unmarshal: options.Unmarshal,
		stop:      make(chan struct{}),
	}
	subscriber.start()
	return subscriber, nil
}

func NewKafkaBroker(opts ...Option) *KafkaBroker {
	options := NewOptions(opts...)
	return &KafkaBroker{
		opts: options,
	}
}
