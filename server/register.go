package server

import (
	"context"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/utils"
	"time"
)

func (g *RPCServer) Register(ctx context.Context) error {

	g.RLock()
	registered := g.registered
	g.RUnlock()

	config := g.opts
	service := g.service.registry

	if !registered {
		log.Infof(ctx, "Registry [%s] Registering node: %s", config.Registry.Name(), service.Name)
	}

	var err error
	for i := 0; i < 3; i++ {
		// attempt to register
		err = config.Registry.Register(ctx, service)
		if err != nil {
			// backoff then retry
			log.Errorf(ctx, "Registry register failed: %s", err.Error())
			time.Sleep(utils.BackoffDelay(i + 1))
			continue
		}
		//log.Debugf("Registry register or keep alive success")
		// success so nil error
		break
	}
	if err != nil {
		log.Errorf(ctx, "Registry [%s] Registering node: %s", config.Registry.Name(), service.Name)
		return err
	}

	if registered {
		return nil
	}

	g.Lock()
	defer g.Unlock()

	if err = g.service.SubscriberAll(ctx); err != nil {
		return err
	}

	g.registered = true
	return nil
}

func (g *RPCServer) Deregister(ctx context.Context) error {

	if g.opts != nil && g.opts.Registry != nil {
		service := g.service.registry
		if err := g.opts.Registry.Deregister(ctx, service); err != nil {
			return err
		}
	}

	g.Lock()
	if !g.registered {
		g.Unlock()
		return nil
	}

	g.registered = false
	for _, err := range g.service.UnsubscribeAll() {
		log.Errorf(ctx, "Unsubscribing from failed: %s", err.Error())
	}
	g.Unlock()
	return nil
}
