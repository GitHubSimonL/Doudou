package logger

import "go.uber.org/zap/zapcore"

//日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

const LogJsonType = "json"
const LogConsoleType = "console"

//默认参数
const (
	Filename          string = "./log/default.log" //日志保存路径 //需要设置程序当前运行路径
	LogLevel          Level  = DebugLevel          //日志记录级别
	MaxSize           int    = 100                 //日志分割的尺寸 MB
	MaxAge            int    = 30                  //分割日志保存的时间 day
	MaxBackups        int    = 1000                //最大日志文件存在个数
	Stacktrace        Level  = PanicLevel          //记录堆栈的级别
	IsCloseStdOut     bool   = false               //是否标准输出console输出
	CallerSkip        int    = 1
	DefaultConsoleLvl        = zapcore.DebugLevel
	MaxAsyncLogSize          = 1024
)

var DefaultOp = Options{
	Filename:      Filename,
	LogLevel:      LogLevel,
	MaxSize:       MaxSize,
	MaxAge:        MaxAge,
	Stacktrace:    Stacktrace,
	IsCloseStdOut: IsCloseStdOut,
	LogType:       LogConsoleType,
	CallerSkip:    CallerSkip,
	MaxBackups:    MaxBackups,
}
