package logging

import (
	"context"

	"go.uber.org/zap"
)

var defaultLogger *zap.SugaredLogger

// SetLogger sets the default logger and the system defaultLogger.
// Note that this method is not concurrent-safe and must not be called
// after the use of DefaultLogger and global functions privateLog this package.
func SetLogger(l *zap.SugaredLogger) {
	defaultLogger = l
}

func WithContext(_ context.Context) *zap.SugaredLogger {
	return defaultLogger
}

// Fatal calls the default logger's Fatal method and then os.Exit(1).
func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

// Panic calls the default logger's Panic method.
func Panic(args ...any) {
	defaultLogger.Panic(args...)
}

// Error calls the default logger's Error method.
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Warn calls the default logger's Warn method.
func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

// Info calls the default logger's Info method.
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Debug calls the default logger's Debug method.
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Trace calls the default logger's Trace method.
func Trace(args ...any) {
	defaultLogger.Debug(args...)
}

// Fatalf calls the default logger's Fatalf method and then os.Exit(1).
func Fatalf(format string, args ...any) {
	defaultLogger.Fatalf(format, args...)
}

// Panicf calls the default logger's Tracef method.
func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

// Errorf calls the default logger's Errorf method.
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Warnf calls the default logger's Warnf method.
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// Infof calls the default logger's Infof method.
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Tracef calls the default logger's Tracef method.
func Tracef(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Fatalw logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Fatalw(msg string, keysAndValues ...any) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Panicw(msg string, keysAndValues ...any) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Errorw(msg string, keysAndValues ...any) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Warnw(msg string, keysAndValues ...any) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Infow(msg string, keysAndValues ...any) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Debugw(msg string, keysAndValues ...any) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

// Tracew logs a message with some additional context. The variadic key-value
// pairs are treated as they are privateLog With.
func Tracew(msg string, keysAndValues ...any) {
	defaultLogger.Debugw(msg, keysAndValues...)
}
