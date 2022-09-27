package zlog

import (
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Format string

const (
	FormatJson    Format = "json"
	FormatConsole Format = "console"

	DefaultFormat     = FormatConsole
	DefaultFileName   = "log/log"
	DefaultMaxFile    = 30
	DefaultCallerSkip = 1
)

type Options struct {
	Format     Format
	FileName   string
	MaxFile    uint
	CallerSkip int
}

var logger *zap.Logger

func init() {
	option := &Options{
		Format:     DefaultFormat,
		FileName:   DefaultFileName,
		MaxFile:    DefaultMaxFile,
		CallerSkip: DefaultCallerSkip,
	}
	SyncFile(option)
}

func SyncFile(option *Options, fields ...zap.Field) {
	logger = zap.New(
		zapcore.NewTee(console(option), file(option)),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(option.CallerSkip),
	)
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

func console(option *Options) zapcore.Core {
	cf := zap.NewProductionEncoderConfig()
	cf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewCore(
		encoder(option.Format, cf),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.DebugLevel,
	)
}

func file(option *Options) zapcore.Core {
	cf := zap.NewProductionEncoderConfig()
	cf.EncodeLevel = zapcore.CapitalLevelEncoder
	wr, err := rotate.New(
		option.FileName+".%Y%m%d.log",
		rotate.WithLinkName(option.FileName),
		rotate.WithRotationCount(option.MaxFile),
		rotate.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.NewCore(
		encoder(option.Format, cf),
		zapcore.AddSync(wr),
		zapcore.InfoLevel,
	)
}

func encoder(f Format, conf zapcore.EncoderConfig) zapcore.Encoder {
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	if f == FormatJson {
		return zapcore.NewJSONEncoder(conf)
	}
	return zapcore.NewConsoleEncoder(conf)
}
