package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

var (
	LOG     *logrus.Entry
	loggers []*Logger
)

func newLogger(name, path string) (*Logger, error) {

	var file *os.File
	if len(path) > 0 {
		var err error
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		file = os.Stdout
	}

	l := &logrus.Logger{
		Out:          file,
		Formatter:    NewFormatter(),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ReportCaller: false,
	}

	logger := &Logger{logger: l, file: file, name: name}
	loggers = append(loggers, logger)
	return logger, nil

}

type Logger struct {
	name   string
	logger *logrus.Logger
	file   *os.File
	lock   sync.Mutex
}

func (l *Logger) Logger() *logrus.Logger {
	return l.logger
}

func (l *Logger) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.file == nil {
		return nil
	}
	path := l.file.Name()
	_, err := os.Stat(l.file.Name())
	if os.IsNotExist(err) {
		file := l.file
		var newFile *os.File
		newFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		l.logger.SetOutput(newFile)
		return file.Close()
	}
	return err
}

func (l *Logger) SetOutput(path string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	var file *os.File
	if l.file != nil {
		file = l.file
	}
	newFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	l.logger.SetOutput(newFile)
	if file != nil {
		return file.Close()
	}
	return nil
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.logger.SetLevel(level)
}

func (l *Logger) SetReportCaller(reportCaller bool) {
	l.logger.SetReportCaller(reportCaller)
}

func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.logger.SetFormatter(formatter)
}

func (l *Logger) AddHook(hook logrus.Hook) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logger.AddHook(hook)
}

func (l *Logger) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.file != nil {
		return nil
	}
	file := l.file
	l.logger.SetOutput(os.Stderr)
	l.file = nil
	return file.Close()
}

func LogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "warning":
		return logrus.WarnLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}
