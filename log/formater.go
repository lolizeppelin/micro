package log

/*
import (
	"bytes"
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	isLinux = runtime.GOOS == "linux"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages func(s string) string

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}
	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	// write message
	if f.TrimMessages != nil {
		b.WriteString(f.TrimMessages(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" %s:%d ",
				entry.Caller.File,
				entry.Caller.Line,
			)
		}
	} else if entry.Level <= logrus.ErrorLevel {
		fmt.Fprintf(b, caller(nil)) // err日志强制打印调用栈
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

func NewFormatter() *Formatter {
	return &Formatter{
		TimestampFormat: utils.TimestampFormat,
		HideKeys:        true,
		NoColors:        isLinux,
		NoFieldsSpace:   true,
		CallerFirst:     true,
		FieldsOrder:     []string{"program"},
	}
}

func ReadableFormater() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	recordTimeFormat := "2006-01-02 15:04:05"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

*/

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

type ConsoleOption func(*ConsoleOptions)

type ConsoleOptions struct {
	// 时间格式 (默认 time.StampMilli)
	TimestampFormat string

	NoColors bool

	Level slog.Leveler
}

func WithTimeFormat(format string) ConsoleOption {
	return func(o *ConsoleOptions) {
		o.TimestampFormat = format
	}
}

func WitOutColor(disabled bool) ConsoleOption {
	return func(o *ConsoleOptions) {
		o.NoColors = disabled
	}
}

func WitHandlerLevel(level slog.Level) ConsoleOption {
	return func(o *ConsoleOptions) {
		o.Level = level
	}
}

func NewConsoleHandler(w io.Writer, opts ...ConsoleOption) *ConsoleHandler {

	options := &ConsoleOptions{
		TimestampFormat: utils.TimestampFormat,
		Level:           slog.LevelInfo,
		NoColors:        true,
	}

	for _, o := range opts {
		o(options)
	}

	return &ConsoleHandler{
		writer:  w,
		options: options,
	}
}

// ConsoleHandler 一般用于控制台调试,方便阅读
type ConsoleHandler struct {
	// 基础配置
	writer  io.Writer
	options *ConsoleOptions

	callerFormatter func(*runtime.Frame) string

	attrs []slog.Attr
	group *group
}

func (h *ConsoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.options.Level != nil {
		minLevel = h.options.Level.Level()
	}
	return level >= minLevel
}

func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	b := &strings.Builder{}

	// 时间戳
	tf := h.options.TimestampFormat
	if tf == "" {
		tf = time.StampMilli
	}
	b.WriteString(r.Time.Format(tf))
	b.WriteByte(' ')

	// 日志级别
	if !h.options.NoColors {
		b.WriteString(h.levelColor(r.Level))
	}
	b.WriteString(strings.ToUpper(r.Level.String()[:4]))
	if !h.options.NoColors {
		b.WriteString("\x1b[0m")
	}

	b.WriteByte(' ')
	stacks := h.writeFields(b, r)
	// 消息处理
	msg := r.Message

	// 调用者信息
	var sc string
	if r.PC != 0 {
		sc = h.formatCaller(r)
	}

	// 组装最终日志行
	var final string

	if sc != "" {
		final = fmt.Sprintf("%s %s%s", b.String(), sc, msg)
	} else {
		final = fmt.Sprintf("%s%s", b.String(), msg)
	}
	_, err := h.writer.Write([]byte(final + "\n"))
	if len(stacks) > 0 {
		for _, attr := range stacks {
			h.writer.Write([]byte(attr.Key + " stack: \n"))
			h.writer.Write([]byte(attr.String()))
		}
	}
	return err
}

func (h *ConsoleHandler) writeFields(b *strings.Builder, r slog.Record) []slog.Attr {
	// 写 本地 入字段
	for _, a := range h.attrs {
		_, _ = fmt.Fprintf(b, "[%s:%v]", a.Key, a.Value.Any())
	}

	var stacks []slog.Attr

	n := r.NumAttrs()
	var attrs []slog.Attr
	if n > 0 {
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})
	}

	if h.group != nil {
		if n == 0 {
			g := h.group.NextNonEmpty()
			if g != nil {
				for _, a := range g.attrs {
					_, _ = fmt.Fprintf(b, "[%s.%s:%v]", g.name, a.Key, a.Value.Any())
				}
			}
		} else {
			g := h.group
			for _, a := range g.attrs {
				_, _ = fmt.Fprintf(b, "[%s.%s:%v]", g.name, a.Key, a.Value.Any())

			}
			for _, a := range attrs {
				if a.Key == "stack" {
					stacks = append(stacks, a)
					continue
				}
				_, _ = fmt.Fprintf(b, "[%s.%s:%v]", g.name, a.Key, a.Value.Any())
			}

			g = g.next
			for g != nil {
				for _, a := range g.attrs {
					_, _ = fmt.Fprintf(b, "[%s.%s:%v]", g.name, a.Key, a.Value.Any())

				}
				for _, a := range attrs {
					_, _ = fmt.Fprintf(b, "[%s.%s:%v]", g.name, a.Key, a.Value.Any())
				}
				g = g.next
			}
		}

	} else if n > 0 { // 没有分组
		// 写 record 入字段
		for _, a := range attrs {
			if a.Key == "stack" {
				stacks = append(stacks, a)
				continue
			}
			_, _ = fmt.Fprintf(b, "[%s:%v]", a.Key, a.Value.Any())
		}
	}

	return stacks
}

func (h *ConsoleHandler) formatCaller(r slog.Record) string {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	frame, _ := fs.Next()

	if h.callerFormatter != nil {
		return h.callerFormatter(&frame)
	}
	return fmt.Sprintf("%s:%d ", frame.File, frame.Line)
}

// 颜色处理
func (h *ConsoleHandler) levelColor(l slog.Level) string {
	if h.options.NoColors {
		return ""
	}
	switch {
	case l >= slog.LevelError:
		return "\x1b[31m" // 红色
	case l >= slog.LevelWarn:
		return "\x1b[33m" // 黄色
	case l >= slog.LevelInfo:
		return "\x1b[36m" // 青色
	default:
		return "\x1b[37m" // 白色
	}
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	if h2.group != nil {
		h2.group = h2.group.Clone()
		h2.group.AddAttrs(attrs)
	} else {
		if h2.attrs == nil {
			h2.attrs = utils.CopySlice(attrs)
		} else {
			h2.attrs = utils.CopySlice(h.attrs)
			h2.attrs = append(h2.attrs, attrs...)
		}
	}
	return &h2
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.group = &group{name: name, next: h2.group}
	return &h2
}

type group struct {
	// name is the name of the group.
	name string
	// attrs are the attributes associated with the group.
	attrs []slog.Attr
	next  *group
}

// NextNonEmpty returns the next group within g's linked-list that has
// attributes (including g itself). If no group is found, nil is returned.
func (g *group) NextNonEmpty() *group {
	if g == nil || len(g.attrs) > 0 {
		return g
	}
	return g.next.NextNonEmpty()
}

// Clone returns a copy of g.
func (g *group) Clone() *group {
	if g == nil {
		return g
	}
	g2 := *g
	g2.attrs = utils.CopySlice(g2.attrs)
	return &g2
}

func (g *group) AddAttrs(attrs []slog.Attr) {
	g.attrs = append(g.attrs, attrs...)
}
