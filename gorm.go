package zlog

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	gl "gorm.io/gorm/logger"
)

type LoggerForGorm struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gl.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type OptionForGorm struct {
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func NewLoggerForGorm(option *OptionForGorm) LoggerForGorm {
	log := LoggerForGorm{
		ZapLogger:                 logger,
		LogLevel:                  gl.Warn,
		SlowThreshold:             option.SlowThreshold,
		SkipCallerLookup:          option.SkipCallerLookup,
		IgnoreRecordNotFoundError: option.IgnoreRecordNotFoundError,
	}
	gl.Default = log
	return log
}

func (l LoggerForGorm) LogMode(level gl.LogLevel) gl.Interface {
	return LoggerForGorm{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l LoggerForGorm) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Info {
		return
	}
	Debugf(str, args...)
}

func (l LoggerForGorm) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Warn {
		return
	}
	Warnf(str, args...)
}

func (l LoggerForGorm) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gl.Error {
		return
	}
	Errorf(str, args...)
}

func (l LoggerForGorm) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gl.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		Errorf("err=%s elapsed=%s rows=%d sql=%s", err.Error(), elapsed.String(), rows, sql)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gl.Warn:
		sql, rows := fc()
		Warnf("err=%s elapsed=%s rows=%d sql=%s", err.Error(), elapsed.String(), rows, sql)
	case l.LogLevel >= gl.Info:
		sql, rows := fc()
		Debugf("err=%s elapsed=%s rows=%d sql=%s", err.Error(), elapsed.String(), rows, sql)
	}
}
