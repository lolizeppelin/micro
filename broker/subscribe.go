package broker

import (
	"context"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"sync"
)

func (k *KafkaBroker) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions(opts...)
	client, err := NewKafkaConsumer(ctx, k.opts.Address, topic, options)
	if err != nil {
		return nil, err
	}

	subscriber := &KafkaSubscriber{
		tracer:    tracing.GetTracer(HandlerScope, _version),
		topic:     topic,
		client:    client,
		handler:   handler,
		fallback:  k.opts.ErrorHandler,
		unmarshal: options.Unmarshal,
		stop:      make(chan struct{}),
		wg:        new(sync.WaitGroup),
	}
	go subscriber.start()
	return subscriber, nil
}

type KafkaSubscriber struct {
	topic     string
	tracer    oteltrace.Tracer
	client    *kgo.Client
	handler   Handler
	stop      chan struct{}
	fallback  ErrorHandler
	wg        *sync.WaitGroup
	unmarshal func([]byte) (*transport.Message, error)
}

func (s *KafkaSubscriber) Topic() string {
	return s.topic
}

func (s *KafkaSubscriber) Unsubscribe() error {
	ctx := context.Background()
	s.client.Close()
	s.stop <- struct{}{}
	s.wg.Wait() // 等待循环推出
	for {       // 处理剩余记录
		fetches := s.client.PollRecords(ctx, 100)
		s.fire(fetches)
		if fetches.IsClientClosed() {
			break
		}
	}
	return nil

}

func (s *KafkaSubscriber) start() {
	ctx := context.Background()
	s.wg.Add(1)
STOP:
	for {
		select {
		case <-s.stop:
			break STOP
		default:
			fetches := s.client.PollRecords(ctx, 100)
			s.fire(fetches)
		}
	}
	s.wg.Done()
}

func (s *KafkaSubscriber) fire(fetches kgo.Fetches) int {
	records := fetches.Records()
	if len(records) <= 0 {
		return 0
	}
	for _, record := range records {
		ctx, msg, err := Decode(record, s.unmarshal)
		if err != nil {
			s.fallback(ctx, "kafka.decode", record, err)
			continue
		}
		event := &kafkaEvent{
			msg: msg,
		}

		var span oteltrace.Span
		ctx, span = s.tracer.Start(ctx, "kafka.consume",
			oteltrace.WithAttributes(
				attribute.String("endpoint", msg.Header[transport.Endpoint]),
			),
		)
		if err = s.handler(ctx, event); err != nil {
			span.RecordError(err)
			s.fallback(ctx, "kafka.handler", record, err)
		}

		span.End()

	}
	return len(records)
}

func Decode(record *kgo.Record, unmarshal ...SubscribeUnmarshal) (ctx context.Context, msg *transport.Message, err error) {
	if len(unmarshal) > 0 {
		msg, err = unmarshal[0](record.Value)
	} else {
		msg, err = _unmarshal(record.Value)
	}

	if err != nil {
		return
	}
	headers := make(map[string]string)
	for _, header := range record.Headers {
		headers[header.Key] = string(header.Value)
	}
	ctx = tracing.Extract(context.Background(), headers)
	return
}
