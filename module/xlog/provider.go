package xlog

type Provider interface {
	Logger() *Logger
}
