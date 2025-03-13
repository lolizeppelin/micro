package broker

import (
	"context"
	"github.com/lolizeppelin/micro/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewKafkaConsumer(address []string, topics string, opts SubscribeOptions) (*kgo.Client, error) {
	options := []kgo.Opt{
		kgo.SeedBrokers(address...),
		kgo.ConsumerGroup(opts.Queue),
		kgo.ConsumeTopics(topics),
		kgo.DisableIdempotentWrite(),
	}
	if opts.AutoAck {
		options = append(options, kgo.RequiredAcks(kgo.NoAck()))
	}
	client, err := kgo.NewClient(options...)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewKafkaProducer(address []string, autoAck ...bool) (*kgo.Client, error) {
	options := []kgo.Opt{
		kgo.SeedBrokers(address...),
		kgo.DisableIdempotentWrite(),
	}
	AutoAck := true
	if len(autoAck) > 0 {
		AutoAck = autoAck[0]
	}
	if AutoAck {
		options = append(options, kgo.RequiredAcks(kgo.NoAck()))
	}
	client, err := kgo.NewClient(options...)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Info("connect kafka success")
	return client, nil
}
