package zlog

import "go.uber.org/zap"

func Debug(msg string, fields ...zap.Field) {
	Logger().Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	Logger().Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	Logger().Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	Logger().Error(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	Logger().Fatal(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	Logger().Panic(msg, fields...)
}
func With(fields ...zap.Field) *zap.Logger {
	return Logger().With(fields...)
}
func WithOptions(fields ...zap.Option) *zap.Logger {
	return Logger().WithOptions(fields...)
}
func Sync() {
	_ = Logger().Sync()
}

func Debugf(template string, args ...interface{}) {
	Sugar().Debugf(template, args...)
}
func Infof(template string, args ...interface{}) {
	Sugar().Infof(template, args...)
}
func Warnf(template string, args ...interface{}) {
	Sugar().Warnf(template, args...)
}
func Errorf(template string, args ...interface{}) {
	Sugar().Errorf(template, args...)
}
func Fatalf(template string, args ...interface{}) {
	Sugar().Fatalf(template, args...)
}
func Panicf(template string, args ...interface{}) {
	Sugar().Panicf(template, args...)
}
func Printf(format string, args ...interface{}) {
	Sugar().Infof(format, args...)
}
func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}
