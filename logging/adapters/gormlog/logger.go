package gormlog

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

var gormSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = sourceDir(file)
}

func sourceDir(file string) string {
	dir := filepath.Dir(filepath.Dir(file))
	s := filepath.Dir(dir)
	if filepath.Base(s) != "gorm.io" {
		s = dir
	}
	return filepath.ToSlash(s) + "/"
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

// GormLogger gorm日志
type GormLogger struct {
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	*zap.SugaredLogger
}

// NewLogger 创建日志
func NewLogger(filename string) *GormLogger {
	gl := &GormLogger{LogLevel: logger.Info, SlowThreshold: 200 * time.Millisecond}
	gl.SugaredLogger = logging.NewLoggerURL("info", filename)
	return gl
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	gl := *l
	gl.LogLevel = level
	return &gl
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		preStr := fmt.Sprintf(infoStr, FileWithLineNum())
		l.SugaredLogger.Infof(preStr+msg, data...)
	}
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		preStr := fmt.Sprintf(warnStr, FileWithLineNum())
		l.SugaredLogger.Warnf(preStr+msg, data...)
	}
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		preStr := fmt.Sprintf(errStr, FileWithLineNum())
		l.SugaredLogger.Errorf(preStr+msg, data...)
	}
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	msecs := float64(elapsed.Nanoseconds()) / 1e6
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Infof(traceErrStr, FileWithLineNum(), err, msecs, "-", sql)
		} else {
			l.Infof(traceErrStr, FileWithLineNum(), err, msecs, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Infof(traceWarnStr, FileWithLineNum(), slowLog, msecs, "-", sql)
		} else {
			l.Infof(traceWarnStr, FileWithLineNum(), slowLog, msecs, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Infof(traceStr, FileWithLineNum(), msecs, "-", sql)
		} else {
			l.Infof(traceStr, FileWithLineNum(), msecs, rows, sql)
		}
	}
}
