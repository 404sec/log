package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Config struct {
	//logger     *zap.Logger
	DebugModel bool
	Level      int
	Encoder    string
	AppName    string
	LogPath    string
	ErrLogPath string
	File       FileConf
	GormOption GormOption
}
type GormOption struct {
	LogLevel                  int
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

/*
Silent LogLevel = iota + 1
Error
Warn
Info
*/
type FileConf struct {
	MaxSize    int  // 每个日志文件保存的大小 单位:M
	MaxAge     int  // 文件最多保存多少天
	MaxBackups int  // 日志文件最多保存多少个备份
	Compress   bool // 是否压缩

}
type Properties struct {
	Core      zapcore.Core
	WritesAll zapcore.WriteSyncer
	WritesErr zapcore.WriteSyncer

	Level zap.AtomicLevel
}

/*
DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
*/
