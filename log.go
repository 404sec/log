package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
)

var globalMu sync.Mutex
var globalLogger, globalProperties, globalSugarLogger atomic.Value

func init() {
	conf := &Config{
		DebugModel: true,
		Level:      0,
		Encoder:    "console",
		AppName:    "app",
		LogPath:    "access.log",
		ErrLogPath: "error.log",
		File: FileConf{
			MaxSize:    128,
			MaxAge:     -1,
			MaxBackups: -1,
			Compress:   false,
		},
		GormOption: GormOption{LogLevel: 3, SlowThreshold: 100 * time.Millisecond, SkipCallerLookup: false, IgnoreRecordNotFoundError: false},
	}
	logger, properties, err := New(conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	ReplaceGlobals(logger, properties)
}
func Init(conf *Config) {

	logger, properties, err := New(conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	ReplaceGlobals(logger, properties)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)
	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)

	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)
	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)
	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)
	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fs := getContext(ctx)
	L().With(fs...).WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// With creates a child logger and adds structured context to it.
// Fields added to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return L().WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// SetLevel alters the logging level.
func SetLevel(l zapcore.Level) {
	globalProperties.Load().(*Properties).Level.SetLevel(l)
}

// GetLevel gets the logging level.
func GetLevel() zapcore.Level {
	return globalProperties.Load().(*Properties).Level.Level()
}

//TODO
func getContext(ctx context.Context) []zap.Field {
	res := make([]zap.Field, 0)
	return res
}
