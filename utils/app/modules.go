package app

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/utils"
)

func startModules(modules []micro.Module) error {
	for _, m := range modules {
		if err := m.Init(); err != nil {
			return err
		}
	}
	for _, m := range modules {
		m.AfterInit()
	}
	return nil
}

func shutdownModules(modules []micro.Module) {
	_modules := utils.SliceReverse(modules)
	for _, m := range _modules {
		m.BeforeShutdown()
	}
	for _, m := range _modules {
		if err := m.Shutdown(); err != nil {
			log.Warnf("error stopping module: %s", err.Error())
		}
	}
	log.Info("all module stopped success")
}
