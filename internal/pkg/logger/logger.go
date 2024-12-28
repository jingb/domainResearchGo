package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 参考文档：
// Zap官方文档: https://pkg.go.dev/go.uber.org/zap
// Zap Github: https://github.com/uber-go/zap
// zapcore.EncoderConfig: https://pkg.go.dev/go.uber.org/zap/zapcore#EncoderConfig

var (
	globalLogger *zap.Logger
	once         sync.Once // 确保 InitLogger 只被执行一次
)

// InitLogger initializes the global logger
func InitLogger() {
	once.Do(func() {
		// EncoderConfig 定义日志的编码方式，决定日志的格式和字段
		// 参考: https://pkg.go.dev/go.uber.org/zap/zapcore#EncoderConfig
		encoderConfig := zapcore.EncoderConfig{
			// 以下字段定义日志中各个部分的键名
			TimeKey:   "time",   // 时间字段的键名
			LevelKey:  "level",  // 日志级别字段的键名
			NameKey:   "logger", // logger名字字段的键名
			CallerKey: "caller", // 调用者信息字段的键名
			//FunctionKey:   zapcore.OmitKey,           // 调用函数信息的键名，OmitKey表示省略该字段
			MessageKey:    "msg",                     // 消息字段的键名
			StacktraceKey: "stacktrace",              // 堆栈跟踪字段的键名
			LineEnding:    zapcore.DefaultLineEnding, // 日志行分隔符，默认是换行符\n

			// 编码器，定义各个字段如何被编码成最终的输出格式
			EncodeLevel: zapcore.CapitalLevelEncoder, // 日志级别编码器，使用大写字母（如 INFO, ERROR）
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05.000")) // 自定义时间格式
			},
			EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码器，将时间转为人类可读格式
			EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用者信息编码器，以包名/文件名:行号 格式显示
		}

		// NewCore创建一个Core，它是日志记录的核心
		// 参考: https://pkg.go.dev/go.uber.org/zap/zapcore#NewCore
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), // 使用普通行格式而不是JSON
			zapcore.AddSync(os.Stdout),               // 输出到标准输出
			zap.NewAtomicLevelAt(zap.InfoLevel),      // 设置日志级别为INFO
		)

		// 创建logger实例
		// 参考: https://pkg.go.dev/go.uber.org/zap#New
		globalLogger = zap.New(core,
			zap.AddCaller(),      // 启用调用者信息记录
			zap.AddCallerSkip(1), // 跳过一层调用栈，确保显示的是调用logger的位置，而不是logger包内部的位置
			zap.AddStacktrace(zap.ErrorLevel),
		)
	})
}

// GetLogger returns the global logger instance
func getLogger() *zap.Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

// Infof logs a formatted message at InfoLevel
func Infof(format string, args ...interface{}) {
	getLogger().Info(fmt.Sprintf(format, args...))
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

// Errorf logs a formatted message at ErrorLevel
func Errorf(format string, args ...interface{}) {
	getLogger().Error(fmt.Sprintf(format, args...))
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

// Debugf logs a formatted message at DebugLevel
func Debugf(format string, args ...interface{}) {
	getLogger().Debug(fmt.Sprintf(format, args...))
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

// Warnf logs a formatted message at WarnLevel
func Warnf(format string, args ...interface{}) {
	getLogger().Warn(fmt.Sprintf(format, args...))
}

// Fatal logs a message at FatalLevel
func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}

// Fatalf logs a formatted message at FatalLevel
func Fatalf(format string, args ...interface{}) {
	getLogger().Fatal(fmt.Sprintf(format, args...))
}
