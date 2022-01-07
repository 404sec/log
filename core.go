package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%04d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/100000))
}
func New(conf *Config) (*zap.Logger, *Properties, error) {

	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	cfg := zapcore.EncoderConfig{
		TimeKey:  "timestamp",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "caller",
		MessageKey: "msg",
		//	StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     formatEncodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if conf.DebugModel {
		cfg.CallerKey = "caller"
		cfg.StacktraceKey = "trace"

	} else {
		cfg.StacktraceKey = "trace"
	}
	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	errLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现

	var writesAll zapcore.WriteSyncer
	var writesErr zapcore.WriteSyncer

	if conf.LogPath != "" {
		infoWriter := WriterHook(conf.LogPath, &conf.File)
		writesAll = zapcore.AddSync(&infoWriter)

	}
	if conf.ErrLogPath != "" {
		errWriter := WriterHook(conf.ErrLogPath, &conf.File)
		writesErr = zapcore.AddSync(&errWriter)
	}
	writesAll = zapcore.AddSync(os.Stdout)

	// 最后创建具体的Logger
	field := zap.Fields(zap.String("appName", conf.AppName))
	var enc zapcore.Encoder
	switch conf.Encoder {
	case "json":
		enc = zapcore.NewJSONEncoder(cfg)
	case "console":
		enc = zapcore.NewConsoleEncoder(cfg)
	case "mapObject":
		zapcore.NewMapObjectEncoder()
	default:
		enc = zapcore.NewJSONEncoder(cfg)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(enc, zapcore.NewMultiWriteSyncer(writesAll), infoLevel),
		zapcore.NewCore(enc, zapcore.NewMultiWriteSyncer(writesErr), errLevel),
	)

	var logger *zap.Logger
	if !conf.DebugModel {
		logger = zap.New(core, zap.AddCaller(), field) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	} else {
		logger = zap.New(core, zap.AddCaller(), field, zap.AddStacktrace(zap.ErrorLevel)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	}

	err := logger.Sync()
	level := zap.NewAtomicLevel()

	level.SetLevel(zapcore.Level(conf.Level))

	r := &Properties{
		Core:      core,
		WritesAll: writesAll,
		WritesErr: writesErr,
		Level:     level,
	}
	return logger, r, err

}
func WriterHook(filename string, config *FileConf) lumberjack.Logger {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	return lumberjack.Logger{
		Filename:   filename,          // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     config.MaxAge,     // 文件最多保存多少天
		MaxBackups: config.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   config.Compress,   // 是否压缩
	}
}
func ReplaceGlobals(logger *zap.Logger, props *Properties) func() {
	globalMu.Lock()
	prevLogger := globalLogger.Load()
	prevProps := globalProperties.Load()
	globalLogger.Store(logger)
	globalSugarLogger.Store(logger.Sugar())
	globalProperties.Store(props)
	globalMu.Unlock()

	if prevLogger == nil || prevProps == nil {
		// When `ReplaceGlobals` is called first time, atomic.Value is empty.
		return func() {}
	}
	return func() {
		ReplaceGlobals(prevLogger.(*zap.Logger), prevProps.(*Properties))
	}
}
func L() *zap.Logger {
	return globalLogger.Load().(*zap.Logger)
}
func Sync() error {
	return L().Sync()
}
