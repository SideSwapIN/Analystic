package logger

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(1))
	if os.Getenv("APP_ENV") == "production" {
		// production environment
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err := cfg.Build(zapOptions...)
		if err != nil {
			panic("failed to initialize logger")
		}
		Logger = logger.Sugar()
	} else {
		// development environment
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err := cfg.Build(zapOptions...)
		if err != nil {
			panic("failed to initialize logger")
		}
		Logger = logger.Sugar()
	}
}

// CloseLogger 关闭日志
func CloseLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func Dump(datas ...interface{}) {
	for i, data := range datas {
		Logger.Debug("Dump-", i)
		bytes, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			Errorf("DumpJson error, print data error: %s", err.Error())
		}
		Logger.Debug(string(bytes))
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	Logger.Debug(GetPrefix(args)...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	Logger.Info(GetPrefix(args)...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	Logger.Warn(GetPrefix(args)...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	Logger.Error(GetPrefix(args)...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// log then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	Logger.DPanic(GetPrefix(args)...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	Logger.Panic(GetPrefix(args)...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	Logger.Fatal(GetPrefix(args)...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	Logger.Debugf(GetPrefixf(template), args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	Logger.Infof(GetPrefixf(template), args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	Logger.Warnf(GetPrefixf(template), args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	Logger.Errorf(GetPrefixf(template), args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// log then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	Logger.DPanicf(GetPrefixf(template), args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	Logger.Panicf(GetPrefixf(template), args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(GetPrefixf(template), args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	Logger.Debugw(GetPrefixf(msg), keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	Logger.Infow(GetPrefixf(msg), keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	Logger.Warnw(GetPrefixf(msg), keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	Logger.Errorw(GetPrefixf(msg), keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// log then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	Logger.DPanicw(GetPrefixf(msg), keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	Logger.Panicw(GetPrefixf(msg), keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	Logger.Fatalw(GetPrefixf(msg), keysAndValues...)
}

func GetPrefix(args []interface{}) []interface{} {
	datas := make([]interface{}, len(args)+1)
	for i, arg := range args {
		datas[i+1] = arg
	}
	return datas
}

func GetPrefixf(template string) string {
	return template
}
