package main

import (
	"context"
	"github.com/404sec/log"
	"go.uber.org/zap"
	"time"
)

func main() {

	conf := &log.Config{
		DebugModel: true,
		Level:      0,
		Encoder:    "json",
		AppName:    "app",
		LogPath:    "access.log",
		ErrLogPath: "error.log",
		File: log.FileConf{
			MaxSize:    128,
			MaxAge:     -1,
			MaxBackups: -1,
			Compress:   false,
		},
		GormOption: log.GormOption{LogLevel: 3, SlowThreshold: 100 * time.Millisecond, SkipCallerLookup: false, IgnoreRecordNotFoundError: false},
	}
	log.Init(conf)

	var t TestTrace
	t.TraceId = "001"
	t.SpanId = "002"
	ctx := context.WithValue(context.TODO(), "jjjj", &t)

	log.Info(ctx, "eee", zap.Int("sss", 111))
}

type TestTrace struct {
	TraceId    string
	SpanId     string
	TraceFlags string
}
