package broker

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"sync"
)

type KafkaSubscriber struct {
	topic     string
	client    *kgo.Client
	handler   Handler
	stop      chan struct{}
	fallback  func(uint8, *kgo.Record, error)
	wg        *sync.WaitGroup
	unmarshal func([]byte) (*transport.Message, error)
}

func (s *KafkaSubscriber) Topic() string {
	return s.topic
}

func (s *KafkaSubscriber) Unsubscribe() error {
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
			s.publish(fetches)
		}
	}
	s.wg.Done()
}

func (s *KafkaSubscriber) publish(fetches kgo.Fetches) int {
	records := fetches.Records()
	if len(records) <= 0 {
		return 0
	}
	for _, record := range records {
		msg, err := s.unmarshal(record.Value)
		if err != nil {
			s.fallback(1, record, err)
			continue
		}
		event := &kafkaEvent{
			msg: msg,
		}
		if err = s.handler(event); err != nil {
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
