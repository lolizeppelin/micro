package log

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/utils"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	handler *multiHandlers
	logger  *Logger
)

func init() {
	opts := new(slog.HandlerOptions)
	opts.Level = slog.LevelInfo
	handler = &multiHandlers{
		opts:     opts,
		files:    utils.NewSyncMap[string, *file](),
		handlers: map[string]slog.Handler{},
		stdout: NewSlog(NewConsoleHandler(os.Stdout,
			WitOutColor(false),
			WitHandlerLevel(opts.Level.Level()))),
		//stdout: NewSlog(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		//	AddSource: true,
		//})),
	}

	logger = NewSlog(handler)
}

type FileHandlerBuilder func(io.Writer, *slog.HandlerOptions) slog.Handler

type file struct {
	path    string
	file    *os.File
	builder FileHandlerBuilder
	sync.RWMutex
}

func (f *file) Write(p []byte) (n int, err error) {
	f.RLock()
	defer f.RUnlock()
	return f.Write(p)
}

func (f *file) Reload() error {
	f.Lock()
	defer f.Unlock()
	if newFile, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return err
	} else {
		oldFile := f.file
		f.file = newFile
		return oldFile.Close()
	}
}

type multiHandlers struct {
	attrs    []slog.Attr
	opts     *slog.HandlerOptions
	files    *utils.SyncMap[string, *file]
	handlers map[string]slog.Handler
	stdout   *Logger
}

func (h *multiHandlers) Level() slog.Level {
	return h.opts.Level.Level()
}

func (h *multiHandlers) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *multiHandlers) Handle(ctx context.Context, record slog.Record) error {
	if len(h.handlers) == 0 {
		return h.stdout.Handler().Handle(ctx, record)
	}
	for name, hdr := range h.handlers {
		if !hdr.Enabled(ctx, record.Level) {
			continue
		}
		if err := hdr.Handle(ctx, record); err != nil {
			h.stdout.ErrorContext(ctx, fmt.Sprintf("logging handler %s record failed: %v", name, err))
		}
	}
	return nil
}

func (h *multiHandlers) WithAttrs(attrs []slog.Attr) slog.Handler {
	copied := &multiHandlers{
		opts:     h.opts,
		handlers: utils.CopyMap(h.handlers),
		stdout:   NewSlog(h.stdout.Handler().WithAttrs(attrs)),
	}
	if len(copied.handlers) == 0 {
		return copied
	}
	for k, v := range copied.handlers {
		h.handlers[k] = v.WithAttrs(attrs)
	}
	return copied
}

func (h *multiHandlers) WithGroup(name string) slog.Handler {

	copied := &multiHandlers{
		opts:     h.opts,
		handlers: utils.CopyMap(h.handlers),
		stdout:   h.stdout.WithGroup(name),
	}
	for k, v := range copied.handlers {
		h.handlers[k] = v.WithGroup(name)
	}
	return copied
}

// AppendFile 添加文件
func (h *multiHandlers) AppendFile(path string, builder ...FileHandlerBuilder) error {
	_, ok := h.files.Load(path)
	if ok {
		return micro.ErrAlreadyExists
	}
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	var f *os.File
	f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	writer := &file{
		path: path,
		file: f,
	}

	var hdr slog.Handler
	if len(builder) > 0 && builder[0] != nil {
		writer.builder = builder[0]
		hdr = writer.builder(writer, h.opts)
	} else {
		hdr = slog.NewTextHandler(writer, h.opts)
	}
	h.files.Store(path, writer)
	h.handlers[path] = hdr.WithAttrs(h.attrs)
	return nil
}

func (h *multiHandlers) Reload() {
	h.files.Range(func(path string, value *file) bool {
		if err := value.Reload(); err != nil {
			h.stdout.Error(fmt.Sprintf("reload file log failed:%v", err))
		}
		return true
	})
}

// AppendHandler 添加非文件的handler
func (h *multiHandlers) AppendHandler(name string, handler slog.Handler) {
	h.handlers[name] = handler.WithAttrs(h.attrs)
}
