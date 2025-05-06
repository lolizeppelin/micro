package log

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/log/internal"
	"log/slog"
	"path/filepath"
)

func Setup(service, program string, logDir string, level slog.Level) error {
	attr := []any{"proc"}
	if program == "" {
		attr = append(attr, service)
	} else {
		attr = append(attr, fmt.Sprintf("%s-%s", service, program))
	}
	if level <= slog.LevelDebug {
		internal.IgnorePC = false
		multi.opts.AddSource = false
	}
	multi.stdout = multi.stdout.With(attr...)
	multi.attrs = argsToAttrSlice(attr)
	multi.opts.Level = level
	if logDir != "" {
		path := filepath.Join(logDir, fmt.Sprintf("%s.log", program))
		if err := multi.AppendFile(path); err != nil {
			return err
		}
	} else {
		multi.handlers["stdout"] = multi.stdout.handler
	}
	return nil
}

func Stack(ctx context.Context, msg string, stack []byte) {
	logger.ErrorContext(ctx, msg, "stack", stack)
}

func Info(ctx context.Context, msg string, attrs ...any) {
	logger.InfoContext(ctx, msg, attrs...)
}

func Warn(ctx context.Context, msg string, attrs ...any) {
	logger.WarnContext(ctx, msg, attrs...)
}

func Warning(ctx context.Context, msg string, attrs ...any) {
	logger.WarnContext(ctx, msg, attrs...)
}

func Error(ctx context.Context, msg string, attrs ...any) {
	logger.ErrorContext(ctx, msg, attrs...)
}

func Trace(ctx context.Context, msg string, attrs ...any) {
	logger.DebugContext(ctx, msg, attrs...)
}

func Debug(ctx context.Context, msg string, attrs ...any) {
	logger.DebugContext(ctx, msg, attrs...)
}

func Fatal(ctx context.Context, msg string, attrs ...any) {
	logger.ErrorContext(ctx, msg, attrs...)
}

func Infof(ctx context.Context, format string, args ...any) {
	logger.InfoContext(ctx, fmt.Sprintf(format, args...))
}

func Warnf(ctx context.Context, format string, args ...any) {
	logger.WarnContext(ctx, fmt.Sprintf(format, args...))
}

func Warningf(ctx context.Context, format string, args ...any) {
	logger.WarnContext(ctx, fmt.Sprintf(format, args...))
}

func Errorf(ctx context.Context, format string, args ...any) {
	logger.ErrorContext(ctx, fmt.Sprintf(format, args...))
}

func Debugf(ctx context.Context, format string, args ...any) {
	logger.DebugContext(ctx, fmt.Sprintf(format, args...))
}

func Tracef(ctx context.Context, format string, args ...any) {
	logger.DebugContext(ctx, fmt.Sprintf(format, args...))
}

func Fatalf(ctx context.Context, format string, args ...any) {
	logger.ErrorContext(ctx, fmt.Sprintf(format, args...))
}

func IsDebugEnabled() bool {
	return slog.LevelDebug <= multi.Level()
}

// AppendHandler 非线程安全, 初始化时调用后不需要在调用
func AppendHandler(name string, h slog.Handler) {
	multi.AppendHandler(name, h)
}

// AppendFileHandler 非线程安全, 初始化时调用后不需要在调用
func AppendFileHandler(name string, builder ...FileHandlerBuilder) error {
	return multi.AppendFile(name, builder...)
}
