package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

type (
	Field    = zapcore.Field
	Provider interface {
		Logger() *Logger
	}
	Logger struct {
		zap   *zap.Logger
		sugar *zap.SugaredLogger
	}
)

var (
	Int64    = zap.Int64
	Int32    = zap.Int32
	Int      = zap.Int
	Uint32   = zap.Uint32
	String   = zap.String
	Any      = zap.Any
	Err      = zap.Error
	NamedErr = zap.NamedError
	Bool     = zap.Bool
	Duration = zap.Duration
)

func New(opts ...Option) *Logger {
	o := newOptions(opts)
	encoder := makeEncoder(o.encoder)
	writer := makeWriter()
	core := zapcore.NewCore(encoder, writer, o.level)
	l := zap.New(core)

	return &Logger{
		zap:   l,
		sugar: l.Sugar(),
	}
}

func makeWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./xlog.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	})
}

func makeEncoder(encoder string) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()

	switch encoder {
	case "json":
		return zapcore.NewJSONEncoder(cfg)
	default:
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewConsoleEncoder(cfg)
	}
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) With(fields ...Field) *Logger {
	ll := *l
	ll.zap = ll.zap.With(fields...)
	return &ll
}

func (l *Logger) Debug(message string, fields ...Field) *Logger {
	l.zap.Debug(message, fields...)
	return l
}

func (l *Logger) Info(message string, fields ...Field) *Logger {
	l.zap.Info(message, fields...)
	return l
}

func (l *Logger) Warn(message string, fields ...Field) *Logger {
	l.zap.Warn(message, fields...)
	return l
}

func (l *Logger) Error(message string, fields ...Field) *Logger {
	l.zap.Error(message, fields...)
	return l
}

func (l *Logger) Fatal(message string, fields ...Field) *Logger {
	l.zap.Fatal(message, fields...)
	return l
}

func (l *Logger) Panic(message string, fields ...Field) *Logger {
	l.zap.Panic(message, fields...)
	return l
}
