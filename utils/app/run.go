package app

import (
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

func Ready() error {
	err := systemd.Ready()
	if errors.Is(err, systemd.NotSupportSystemd) {
		return nil
	}
	return err
}

func Stopping() error {
	err := systemd.Stopping()
	if errors.Is(err, systemd.NotSupportSystemd) {
		return nil
	}
	return err
}

func Run(modules []micro.Module, daemon bool) error {
	sg := make(chan os.Signal)

	defer shutdownModules(modules)
	err := startModules(modules)
	if err != nil {
		return err
	}

	if daemon {
		if err = Ready(); err != nil {
			return err
		}
		signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT)
		select {
		case <-die:
			log.Info("the app will shutdown in a few seconds")
		case s := <-sg:
			log.Info("------------------------------------------")
			log.Infof("got signal: '%v',  shutting down...", s)
			log.Info("------------------------------------------")
			close(die)
		}
		return Stopping()
	}
	return nil
}

func Stop() {
	close(die)
}
