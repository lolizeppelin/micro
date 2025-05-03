package app

import (
	"context"
	"errors"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/utils/systemd"
	"os"
	"os/signal"
	"syscall"
)

var (
	die = make(chan struct{})
)

func ready() error {
	err := systemd.Ready()
	if errors.Is(err, systemd.NotSupportSystemd) {
		return nil
	}
	return err
}

func stopping() error {
	err := systemd.Stopping()
	if errors.Is(err, systemd.NotSupportSystemd) {
		return nil
	}
	return err
}

func Run(ctx context.Context, modules []micro.Module, daemon bool) error {
	sg := make(chan os.Signal)

	defer shutdownModules(ctx, modules)
	err := startModules(modules)
	if err != nil {
		return err
	}

	if daemon {
		if err = ready(); err != nil {
			return err
		}
		signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT)
		select {
		case <-die:
			log.Info(ctx, "the app will shutdown in a few seconds")
		case s := <-sg:
			log.Info(ctx, "------------------------------------------")
			log.Infof(ctx, "got signal: '%v',  shutting down...", s)
			log.Info(ctx, "------------------------------------------")
			close(die)
		}
		return stopping()
	}
	return nil
}

func Stop() {
	close(die)
}
