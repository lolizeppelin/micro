package broker

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/vmihailenco/msgpack/v5"
	"sync"
)

type kafkaSubscriber struct {
	topic    string
	client   *kgo.Client
	handler  Handler
	stop     chan struct{}
	fallback func(uint8, *kgo.Record, error)
	wg       *sync.WaitGroup
}

func (s *kafkaSubscriber) Topic() string {
	return s.topic
}

func (s *kafkaSubscriber) Unsubscribe() error {
	ctx := context.Background()
	s.client.Close()
	s.wg.Wait() // 等待循环推出
	for {       // 处理剩余记录
		fetches := s.client.PollRecords(ctx, 100)
		if s.publish(fetches) <= 0 {
			break
		}
	}
	return nil

}

func (s *kafkaSubscriber) start() {
	ctx := context.Background()
	s.wg.Add(1)
STOP:
	for {
		select {
		case <-s.stop:
			break STOP
		default:
			fetches := s.client.PollRecords(ctx, 100)
			s.publish(fetches)
		}
	}
	s.wg.Done()
}

func (s *kafkaSubscriber) publish(fetches kgo.Fetches) int {
	records := fetches.Records()
	if len(records) <= 0 {
		return 0
	}
	for _, record := range records {
		// resource := string(record.Key)
		msg := new(transport.Message)
		if err := msgpack.Unmarshal(record.Value, msg); err != nil {
			s.fallback(1, record, err)
			continue
		}
		event := &kafkaEvent{
			msg: msg,
		}
		if err := s.handler(event); err != nil {
			s.fallback(2, record, err)
		}

	}
	return len(records)
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
