package broker

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/vmihailenco/msgpack/v5"
)

func (k *KafkaBroker) Publish(ctx context.Context, topic string, msg *transport.Message) error {
	if k.producer == nil {
		return fmt.Errorf("producer not connect")
	}
	buff, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	headers := make([]kgo.RecordHeader, 0)
	carrier := make(map[string]string)
	tracing.Inject(ctx, carrier)

	for key, value := range carrier {
		headers = append(headers, kgo.RecordHeader{
			Key:   key,
			Value: []byte(value),
		})
	}

	record := &kgo.Record{
		Topic:   topic,
		Value:   buff,
		Headers: headers,
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
