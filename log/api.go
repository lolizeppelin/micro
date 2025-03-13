package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func Setup(program string, logDir string, level logrus.Level) error {
	if LOG != nil {
		panic("log already setup")
	}
	l, err := Mount(logDir, fmt.Sprintf("%s.log", program), level)
	if err != nil {
		return err
	}
	formatter := NewFormatter()
	if level >= logrus.DebugLevel {
		formatter.CustomCallerFormatter = caller
		l.SetReportCaller(true)
	}
	l.SetFormatter(formatter)
	LOG = l.Logger().WithFields(logrus.Fields{"program": program})
	// 发送SIGTERM给自身进程
	l.Logger().ExitFunc = func(i int) {
		p, _ := os.FindProcess(os.Getpid())
		if p != nil {
			LOG.Warnf("logrus send exit signal")
			_ = p.Signal(syscall.SIGTERM)
		}
	}
	return nil
}

func SetMetric(m LoggingMetric) {
	metric = m
}

func Info(args ...interface{}) {
	LOG.Info(args...)
}

func Warn(args ...interface{}) {
	metric.Warn()
	LOG.Warn(args...)
}

func Warning(args ...interface{}) {
	metric.Warn()
	LOG.Warn(args...)
}

func Error(args ...interface{}) {
	metric.Error()
	LOG.Error(args...)
}

func Trace(args ...interface{}) {
	LOG.Trace(args...)
}

func Debug(args ...interface{}) {
	LOG.Debug(args...)
}

func Fatal(args ...interface{}) {
	LOG.Fatal(args...)
}

func Infof(format string, args ...interface{}) {
	LOG.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	metric.Warn()
	LOG.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	metric.Warn()
	LOG.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	metric.Error()
	LOG.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	LOG.Debugf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	LOG.Tracef(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	LOG.Fatalf(format, args...)
}

func IsDebugEnabled() bool {
	return LOG.Logger.IsLevelEnabled(logrus.DebugLevel)
}
