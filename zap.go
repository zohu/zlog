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

var logger *zap.Logger

func init() {
	SyncFile(&Option{
		Format:     DefaultFormat,
		FileName:   DefaultFileName,
		MaxFile:    DefaultMaxFile,
		CallerSkip: DefaultCallerSkip,
	})
}

type Option struct {
	Format     Format
	FileName   string
	MaxFile    uint
	CallerSkip int
}

func SyncFile(ops *Option) {
	logger = zap.New(
		zapcore.NewTee(console(ops), file(ops)),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(ops.CallerSkip),
	)
	zap.ReplaceGlobals(logger)
}

func Logger() *zap.Logger {
	return logger
}

func console(opt *Option) zapcore.Core {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewCore(
		encoder(opt.Format, conf),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.DebugLevel,
	)
}

func file(opt *Option) zapcore.Core {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	wr, err := rotate.New(
		opt.FileName+".%Y%m%d.log",
		rotate.WithLinkName(opt.FileName),
		rotate.WithRotationCount(opt.MaxFile),
		rotate.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.NewCore(
		encoder(opt.Format, conf),
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
