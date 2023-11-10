package xormlog

import (
	"fmt"

	"github.com/azhai/gozzo/logging"
	"go.uber.org/zap"
	"xorm.io/xorm/log"
)

// XormLogger xorm日志
type XormLogger struct {
	level   log.LogLevel
	showSQL bool
	*zap.SugaredLogger
}

// NewLogger 创建日志
func NewLogger(filename string) *XormLogger {
	l := logging.NewLoggerURL("info", filename)
	return WrapLogger(l)
}

// WrapLogger 封装日志
func WrapLogger(l *zap.SugaredLogger) *XormLogger {
	lvl := log.LOG_INFO
	if l == nil {
		l, lvl = zap.NewNop().Sugar(), log.LOG_OFF
	}
	return &XormLogger{level: lvl, showSQL: true, SugaredLogger: l}
}

// AfterSQL implements ContextLogger
func (l *XormLogger) AfterSQL(ctx log.LogContext) {
	var sessionPart string
	v := ctx.Ctx.Value(log.SessionIDKey)
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
	}
	if ctx.ExecuteTime > 0 {
		l.Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

// BeforeSQL implements ContextLogger
func (l *XormLogger) BeforeSQL(log.LogContext) {
}

// Level implement log.Logger
func (l *XormLogger) Level() log.LogLevel {
	return l.level
}

// SetLevel implement log.Logger
func (l *XormLogger) SetLevel(level log.LogLevel) {
	l.level = level
	return
}

// ShowSQL implement log.Logger
func (l *XormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

// IsShowSQL implement log.Logger
func (l *XormLogger) IsShowSQL() bool {
	return l.showSQL
}
