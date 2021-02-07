package logx

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

var traceContent = []string{
	"level", "@timestamp", "trace", "duration", sourceDir, "content",
}

const sourceDir = "source"

type traceLogger struct {
	logEntry
	Trace string `json:"trace,omitempty"`
	Span  string `json:"span,omitempty"`
	ctx   context.Context
}

func (l *traceLogger) Error(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprint(v...), durationCallerDepth))
	}
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprintf(format, v...), durationCallerDepth))
	}
}

func (l *traceLogger) Info(v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Slow(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *traceLogger) write(writer io.Writer, level, content string) {
	l.Timestamp = getTimestamp()
	l.Level = level
	l.Content = content
	l.Trace = traceIdFromContext(l.ctx)
	l.Span = spanIdFromContext(l.ctx)
	outputTrace(writer, l)
}

func outputTrace(writer io.Writer, l *traceLogger) {
	buf := strings.Builder{}
	mp := Convert2Map(l)

	for _, k := range traceContent {
		if k == sourceDir {
			buf.WriteString(FileWithLineNum())
			buf.WriteString(" ")
			continue
		}

		if v, ok := mp[k]; ok {
			if s, ok := v.(string); ok {
				buf.WriteString(s)
			} else {
				buf.WriteString(JsonMarshal(v))
			}
			buf.WriteString(" ")
		} else {
			buf.WriteString("- ")
		}
	}

	s := buf.String()
	writer.Write([]byte(s))
}

func WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
	}
}

func spanIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(tracespec.TracingKey).(tracespec.Trace)
	if !ok {
		return ""
	}

	return t.SpanId()
}

func traceIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(tracespec.TracingKey).(tracespec.Trace)
	if !ok {
		return ""
	}

	return t.TraceId()
}
