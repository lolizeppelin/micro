//go:build !aix && !windows
// +build !aix,!windows

package log

import (
	"os"
	"os/signal"
	"syscall"
)

var c chan os.Signal

func init() { //监控退出信号关闭日志, 监控SIGUSR1 重新打开日志
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for {
			_, ok := <-c
			if !ok {
				break
			}
			for _, logger := range loggers {
				logger.Reload()
			}
		}
	}()
}

func Close() {
	signal.Stop(c)
	close(c)
	for _, logger := range loggers {
		logger.Close()
	}
}
