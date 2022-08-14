package logger

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func LoggerInit() {
	ll := lumberjack.Logger{
		Filename:   "/tmp/logs",
		MaxSize:    1, //MB
		MaxBackups: 3,
		MaxAge:     1, //days
		Compress:   true,
	}
	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &ll,
		}, nil
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "call",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:      "console",
		EncoderConfig: encoderConfig,
		OutputPaths:   []string{fmt.Sprintf("lumberjack:%s", "/tmp/logs")},
	}

	var err error
	Logger, err = loggerConfig.Build()
	if err != nil {
		panic(err)
	}
}

func Trace(level string, lang string, msg string) bool {
	defer Logger.Sync()
	sugar := Logger.Sugar()

	sugar.Infow(msg,
		"level", level,
		"lang", level,
	)

	return true
}
