package broker

import (
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaBroker struct {
	producer *kgo.Client
	opts     *Options
}

func (k *KafkaBroker) Name() string {
	return "kafka"
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

func NewKafkaBroker(opts ...Option) *KafkaBroker {
	options := NewOptions(opts...)
	return &KafkaBroker{
		opts: options,
	}
}
