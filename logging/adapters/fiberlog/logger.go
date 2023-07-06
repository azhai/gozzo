package fiberlog

import (
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
	fl := &FiberLogger{level: log.LevelTrace}
	fl.SugaredLogger = logging.NewLoggerURL("debug", filename)
	return fl
}

// Trace log a message
func (l *FiberLogger) Trace(v ...any) {
	l.SugaredLogger.Debug(v...)
}

// Tracef format and log a message
func (l *FiberLogger) Tracef(format string, v ...any) {
	l.SugaredLogger.Debugf(format, v...)
}

// Tracew log a message of some pairs
func (l *FiberLogger) Tracew(msg string, keysAndValues ...any) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}
