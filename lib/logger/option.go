package logger

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type Option func(*Options)

type Options struct {
	Filename      string // 日志保存路径
	LogLevel      Level  // 日志记录级别
	MaxSize       int    // 日志分割的尺寸 MB
	MaxAge        int    // 分割日志保存的时间 day
	Stacktrace    Level  // 记录堆栈的级别
	IsCloseStdOut bool   // 是否关闭标准输出console输出
	LogType       string // 日志类型,普通 或 json
	MaxBackups    int    // 最大日志文件存在数(n+1)
	IsASync       bool   // 是否异步写
	IsCompress    bool   // 是否压缩
	IsOpenPprof   bool   // 是否打开pprof
	CallerSkip    int
}

func WithLogType(logType string) Option {
	return func(o *Options) {
		o.LogType = logType
	}
}

func WithFilename(logPath string) Option {
	return func(o *Options) {
		o.Filename = logPath
	}
}

func WithLogLevel(logLevel Level) Option {
	return func(o *Options) {
		o.LogLevel = logLevel
	}
}

func WithMaxSize(maxSize int) Option {
	return func(o *Options) {
		o.MaxSize = maxSize
	}
}

func WithMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.MaxAge = maxAge
	}
}

func WithStacktrace(stacktrace Level) Option {
	return func(o *Options) {
		o.Stacktrace = stacktrace
	}
}

// 是否关闭标准输出
func WithIsCloseStdOut(isCloseStdout bool) Option {
	return func(o *Options) {
		o.IsCloseStdOut = isCloseStdout
	}
}

func WithCallerSkip(callerSkip int) Option {
	return func(o *Options) {
		o.CallerSkip = callerSkip
	}
}

func WithBackups(backups int) Option {
	return func(o *Options) {
		o.MaxBackups = backups
	}
}

func WithIsASync(isASync bool) Option {
	return func(o *Options) {
		o.IsASync = isASync
	}
}

func WithIsCompress(isCompress bool) Option {
	return func(o *Options) {
		o.IsCompress = isCompress
	}
}

func WithPprof(isOpenPprof bool) Option {
	return func(o *Options) {
		o.IsOpenPprof = isOpenPprof
	}
}

func (ops *Options) parseLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	}

	return zapcore.DebugLevel
}

func (ops *Options) getLogger() (lowLogger, highLogger, esLogger *LWriter) {
	lowLogger = NewLWriter(ops.Filename, ops.MaxSize, ops.MaxBackups, ops.MaxAge, ops.IsCompress, ops.IsASync)
	highLogger = NewLWriter(ops.Filename+".err", ops.MaxSize, ops.MaxBackups, ops.MaxAge, ops.IsCompress, false)

	// 级别较高的日志默认不允许异步写
	esLogger = NewLWriter(strings.Replace(ops.Filename, ".log", ".es.log", 1), ops.MaxSize, ops.MaxBackups, ops.MaxAge, ops.IsCompress, false) // esLog落地
	return
}
