package app

import (
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

func Run(modules []micro.Module, daemon bool) error {
	sg := make(chan os.Signal)

	defer shutdownModules(modules)
	err := startModules(modules)
	if err != nil {
		return err
	}

	if daemon {
		if err = systemd.Ready(); err != nil {
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
		return systemd.Stopping()
	}
	return nil
}

func Stop() {
	close(die)
}
