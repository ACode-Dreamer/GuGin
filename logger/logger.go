package logger

import (
	"go.uber.org/zap"
	"sync"
)

// MyLogger 是一个自定义的日志记录器结构体
type MyLogger struct {
	*zap.SugaredLogger
}

var myLogger *MyLogger
var loggerOnce sync.Once

// init 初始化日志记录器
func InitLogger() {
	// 创建 SugaredLogger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	// 设置输出到终端，并启用颜色
	_ = zap.ReplaceGlobals(logger)

	myLogger = &MyLogger{sugar}
}

func getLogger() *MyLogger {
	loggerOnce.Do(func() {
		InitLogger()
	})
	return myLogger
}

// Debug 记录 Debug 级别的日志
func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}

// Info 记录 Info 级别的日志
func Info(args ...interface{}) {
	getLogger().Info(args...)
}

// Warn 记录 Warn 级别的日志
func Warn(args ...interface{}) {
	getLogger().Warn(args...)
}

// Error 记录 Error 级别的日志
func Error(args ...interface{}) {
	getLogger().Error(args...)
}

// Fatal 记录 Fatal 级别的日志
func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}

// Panic 记录 Panic 级别的日志
func Panic(args ...interface{}) {
	getLogger().Panic(args...)
}
