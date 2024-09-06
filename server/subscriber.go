package server

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
	"github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/registry"
	"github.com/lolizeppelin/micro/transport"
	"runtime/debug"
	"strings"
	"sync"
)

// dispatch 处理broker数据
func (s *Service) dispatch(event broker.Event) (err error) {
	wg := s.opts.WaitGroup

	wg.Add(1)

	defer func() {
		wg.Done()
		if e := event.Ack(); e != nil {
			log.Errorf("brcker ack failed： %s", err.Error())
		}
		if r := recover(); r != nil {
			log.Errorf("panic recovered: \n%s", string(debug.Stack()))
			err = errors.InternalServerError("go.micro.server", "panic recovered: %v", r)
		}
	}()

	msg := event.Message()
	if msg.Header == nil {
		return fmt.Errorf("no haaders found")
	}

	endpoint := msg.Header[transport.Endpoint]
	parts := strings.Split(endpoint, ".")
	if len(parts) != 2 {
		return fmt.Errorf("no haaders found")
	}
	service := parts[0]
	method := parts[1]
	services, ok := s.services[service]
	if !ok {
		return fmt.Errorf("no haaders found")
	}
	handler := services[method]
	if handler == nil {
		return fmt.Errorf("handler found")
	}
	if handler.Response != nil {
		return fmt.Errorf("not a subscriber handler")
	}
	if handler.Metadata["req"] == "stream" {
		return fmt.Errorf("not support subscriber for stream")
	}

	hdr := make(map[string]string, len(msg.Header))
	for k, v := range msg.Header {
		if k == micro.ContentType {
			continue
		}
		hdr[k] = v
	}

	ctx := transport.NewContext(context.Background(), hdr)
	args, err := handler.BuildArgs(ctx, msg.Header[micro.ContentType], msg.Query, msg.Body)
	if err != nil {
		return err
	}
	handler.Method.Func.Call(args)
	return nil
}

func (s *Service) SubscriberAll() error {
	if s.opts.Broker == nil {
		return nil
	}

	topics := registry.Topics(s.opts.Name, s.opts.Version, s.registry.Nodes[0].Id)

	for _, topic := range topics {
		if _, ok := s.subscribed[topic]; ok {
			continue
		}
		sub, e := s.opts.Broker.Subscribe(topic, s.dispatch, s.opts.BrokerOpts...)

		if e != nil {
			return e
		}
		s.subscribed[topic] = sub
	}

	return nil
}

func (s *Service) UnsubscribeAll() (errors []error) {
	wg := new(sync.WaitGroup)
	for topic, sub := range s.subscribed {

		wg.Add(1)
		go func(s broker.Subscriber) {
			defer wg.Done()
			if err := s.Unsubscribe(); err != nil {
				errors = append(errors, fmt.Errorf("unsubscribe topic %s failed: %s", s.Topic(), err.Error()))
			}
		}(sub)

		delete(s.subscribed, topic)
	}
	wg.Wait()

	return
}
