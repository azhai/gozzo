package gormlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/azhai/gozzo/logging"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

// GormLogger gorm日志
type GormLogger struct {
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	*zap.SugaredLogger
}

// NewLogger 创建日志
func NewLogger(filename string) *GormLogger {
	l := logging.NewLoggerURL("info", filename)
	return WrapLogger(l)
}

// WrapLogger 封装日志
func WrapLogger(l *zap.SugaredLogger) *GormLogger {
	lvl := logger.Info
	if l == nil {
		l, lvl = zap.NewNop().Sugar(), logger.Silent
	}
	s := 200 * time.Millisecond
	return &GormLogger{LogLevel: lvl, SlowThreshold: s, SugaredLogger: l}
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	gl := *l
	gl.LogLevel = level
	return &gl
}

// Info print info
func (l *GormLogger) Info(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		preStr := fmt.Sprintf(infoStr, logging.FileWithLineNum())
		l.SugaredLogger.Infof(preStr+msg, data...)
	}
}

// Warn print warn messages
func (l *GormLogger) Warn(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		preStr := fmt.Sprintf(warnStr, logging.FileWithLineNum())
		l.SugaredLogger.Warnf(preStr+msg, data...)
	}
}

// Error print error messages
func (l *GormLogger) Error(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		preStr := fmt.Sprintf(errStr, logging.FileWithLineNum())
		l.SugaredLogger.Errorf(preStr+msg, data...)
	}
}

// Trace print sql message
func (l *GormLogger) Trace(_ context.Context, begin time.Time,
	fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	microSec := float64(elapsed.Nanoseconds()) / 1e6
	lineNo := logging.FileWithLineNum()
	sql, rows := fc()
	switch {
	case err != nil && l.LogLevel >= logger.Error && !l.IsIgnoreNotFound(err):
		if rows == -1 {
			l.Infof(traceErrStr, lineNo, err, microSec, "-", sql)
		} else {
			l.Infof(traceErrStr, lineNo, err, microSec, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Infof(traceWarnStr, lineNo, slowLog, microSec, "-", sql)
		} else {
			l.Infof(traceWarnStr, lineNo, slowLog, microSec, rows, sql)
		}
	case l.LogLevel == logger.Info:
		if rows == -1 {
			l.Infof(traceStr, lineNo, microSec, "-", sql)
		} else {
			l.Infof(traceStr, lineNo, microSec, rows, sql)
		}
	}
}

// IsIgnoreNotFound when we want to ignore NotFound Record error
func (l *GormLogger) IsIgnoreNotFound(err error) bool {
	return l.IgnoreRecordNotFoundError && errors.Is(err, logger.ErrRecordNotFound)
}
