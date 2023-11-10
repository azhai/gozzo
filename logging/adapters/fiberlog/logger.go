package fiberlog

import (
	"context"
	"io"

	"github.com/azhai/gozzo/logging"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
)

// FiberLogger fiber日志
type FiberLogger struct {
	level log.Level
	*zap.SugaredLogger
}

// NewLogger 创建日志
func NewLogger(filename string) *FiberLogger {
	l := logging.NewLoggerURL("info", filename)
	return WrapLogger(l)
}

// WrapLogger 封装日志
func WrapLogger(l *zap.SugaredLogger) *FiberLogger {
	lvl := log.LevelInfo
	if l == nil {
		l, lvl = zap.NewNop().Sugar(), log.LevelPanic
	}
	return &FiberLogger{level: lvl, SugaredLogger: l}
}

func (l *FiberLogger) Trace(v ...interface{}) {
	l.Debug(v...)
}

func (l *FiberLogger) Tracef(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func (l *FiberLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.Debugw(msg, keysAndValues...)
}

// SetLevel implement log.ControlLogger
func (l *FiberLogger) SetLevel(level log.Level) {
	l.level = level
}

// SetOutput implement log.ControlLogger
func (l *FiberLogger) SetOutput(_ io.Writer) {
}

// WithContext implement log.AllLogger
func (l *FiberLogger) WithContext(_ context.Context) log.CommonLogger {
	return l
}
