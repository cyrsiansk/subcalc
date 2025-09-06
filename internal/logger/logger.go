package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(levelStr string, useFileRotation bool) (*zap.Logger, error) {
	var level zapcore.Level
	switch levelStr {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.TimeKey = "ts"
	enc := zapcore.NewJSONEncoder(encCfg)

	var ws zapcore.WriteSyncer
	if useFileRotation {
		l := &lumberjack.Logger{
			Filename:   "/var/log/subscriptions/app.log",
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     14,
			Compress:   true,
		}
		ws = zapcore.AddSync(l)
	} else {
		ws = zapcore.Lock(os.Stdout)
	}

	core := zapcore.NewCore(enc, ws, zap.NewAtomicLevelAt(level))

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	logger = logger.With(
		zap.String("service", "subscriptions"),
		zap.String("env", os.Getenv("APP_ENV")),
		zap.Time("start_time", time.Now().UTC()),
	)

	return logger, nil
}
