package zlog

import (
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	DefaultFormat       = Format_CONSOLE
	DefaultFileName     = "log/log"
	DefaultMaxFile      = 30
	DefaultCallerEnable = false
	DefaultCallerSkip   = 1
)

var logger *zap.Logger

func init() {
	option := &Config{
		Format:       DefaultFormat,
		FileName:     DefaultFileName,
		MaxFile:      DefaultMaxFile,
		CallerEnable: DefaultCallerEnable,
		CallerSkip:   DefaultCallerSkip,
	}
	SyncFile(option)
}

func SyncFile(conf *Config, fields ...zap.Field) {
	logger = zap.New(
		zapcore.NewTee(console(conf), file(conf)),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if conf.GetCallerEnable() {
		logger = logger.WithOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(int(conf.GetCallerSkip())),
		)
	}
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	zap.ReplaceGlobals(logger)
}

func Logger() *zap.Logger {
	return logger
}

func ReplaceGlobals(l *zap.Logger) {
	logger = l
	zap.ReplaceGlobals(l)
}

func console(conf *Config) zapcore.Core {
	cf := zap.NewProductionEncoderConfig()
	if conf.GetFormat() == Format_JSON {
		cf.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		cf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return zapcore.NewCore(
		encoder(conf.GetFormat(), cf),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.DebugLevel,
	)
}

func file(conf *Config) zapcore.Core {
	cf := zap.NewProductionEncoderConfig()
	cf.EncodeLevel = zapcore.CapitalLevelEncoder
	wr, err := rotate.New(
		conf.GetFileName()+".%Y%m%d.log",
		rotate.WithLinkName(conf.GetFileName()),
		rotate.WithRotationCount(uint(conf.GetMaxFile())),
		rotate.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.NewCore(
		encoder(conf.GetFormat(), cf),
		zapcore.AddSync(wr),
		zapcore.InfoLevel,
	)
}

func encoder(f Format, conf zapcore.EncoderConfig) zapcore.Encoder {
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	if f == Format_JSON {
		return zapcore.NewJSONEncoder(conf)
	}
	return zapcore.NewConsoleEncoder(conf)
}
