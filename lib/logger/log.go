package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var instance *Logger
var once sync.Once

type Logger struct {
	file     *zap.SugaredLogger
	console  *zap.SugaredLogger //默认初始化console
	esLog    *zap.SugaredLogger
	atom     zap.AtomicLevel
	fileName string

	ops *Options
}

func init() {
	instance = new(Logger)
	instance.atom = zap.NewAtomicLevel()
	instance.atom.SetLevel(DefaultConsoleLvl)

	writeConsole := zapcore.AddSync(os.Stdout)
	allPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= instance.atom.Level()
	})

	enCfg := GetEncoder(LogConsoleType)
	core := zapcore.NewTee(
		zapcore.NewCore(enCfg, writeConsole, allPriority),
	)

	logger := GetZapLogger(core, CallerSkip)
	if logger == nil {
		panic("logger init fail.")
	}

	instance.console = logger
}

func GetEncoder(typ string) zapcore.Encoder {
	enConfig := zap.NewProductionEncoderConfig()
	enConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02-15:04:05")

	var enCfg zapcore.Encoder = nil
	switch typ {
	case LogJsonType:
		enCfg = zapcore.NewJSONEncoder(enConfig)
	case LogConsoleType:
		enConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		fallthrough
	default:
		enCfg = zapcore.NewConsoleEncoder(enConfig)
	}

	return enCfg
}

func GetZapLogger(core zapcore.Core, callerSkip int) *zap.SugaredLogger {
	opList := make([]zap.Option, 0)
	opList = append(opList, zap.AddStacktrace(zapcore.PanicLevel))
	opList = append(opList, zap.AddCaller())
	opList = append(opList, zap.AddCallerSkip(callerSkip))

	zapLogger := zap.New(core, opList...).Sugar()
	if zapLogger == nil {
		panic("zaplooger get fail.")
		return nil
	}

	return zapLogger
}

func InitLogger(opts ...Option) {
	once.Do(func() {
		instance.Init(opts...)
	})

	err := redirectStdErrLog()
	if err != nil {
		LogErrf("Failed to redirect stderr to file: %v", err)
	}
}

func (log *Logger) Init(opts ...Option) {
	o := DefaultOp

	for _, op := range opts {
		op(&o)
	}

	log.ops = &o

	lowLog, highLog, esLog := o.getLogger()

	if lowLog == nil || highLog == nil {
		panic("log init fail.")
		return
	}

	log.atom.SetLevel(o.parseLevel(o.LogLevel))

	// 打印到控制台和文件
	writeFile1 := zapcore.AddSync(lowLog)
	writeFile2 := zapcore.AddSync(highLog)
	esFile := zapcore.AddSync(esLog)

	lowWS := zapcore.NewMultiWriteSyncer(writeFile1)
	highWS := zapcore.NewMultiWriteSyncer(writeFile2)
	esWS := zapcore.NewMultiWriteSyncer(esFile)

	if o.IsCloseStdOut {
		log.console = nil
	}

	if o.IsOpenPprof {
		StartPprofWork()
	}

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= log.atom.Level()
	})

	esPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev <= zap.ErrorLevel && lev >= log.atom.Level()
	})

	//生成配置
	enCfg := GetEncoder(o.LogType)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(enCfg, lowWS, lowPriority),
		zapcore.NewCore(enCfg, highWS, highPriority),
	)

	esCore := zapcore.NewTee(
		zapcore.NewCore(GetEncoder(LogJsonType), esWS, esPriority),
	)

	logger := GetZapLogger(core, o.CallerSkip)
	if logger == nil {
		panic("logger is nil.")
	}

	log.file = logger
	log.esLog = GetZapLogger(esCore, o.CallerSkip)
	log.fileName = o.Filename
}

func (log *Logger) Close() error {
	if log.file == nil {
		return nil
	}

	return log.file.Sync()
}

func Log(args ...interface{}) {
	if instance.console != nil {
		instance.console.Info(args...)
	}

	if instance.file != nil {
		instance.file.Info(args...)
	}
}

func Logf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Infof(template, args...)
	}

	if instance.file != nil {
		instance.file.Infof(template, args...)
	}
}

func LogDebug(args ...interface{}) {
	if instance.console != nil {
		instance.console.Debug(args)
	}

	if instance.file != nil {
		instance.file.Debug(args)
	}
}

func LogInfo(args ...interface{}) {
	if instance.console != nil {
		instance.console.Info(args...)
	}

	if instance.file != nil {
		instance.file.Info(args...)
	}
}

func LogPrint(args ...interface{}) {
	if instance.console != nil {
		instance.console.Info(args)
	}

	if instance.file != nil {
		instance.file.Info(args)
	}
}

func LogWarn(args ...interface{}) {
	if instance.console != nil {
		instance.console.Warn(args)
	}

	if instance.file != nil {
		instance.file.Warn(args)
	}
}

func LogErr(args ...interface{}) {
	if instance.console != nil {
		instance.console.Error(args)
	}

	if instance.file != nil {
		instance.file.Error(args)
	}
}

func LogDPanic(args ...interface{}) {
	if instance.console != nil {
		instance.console.DPanic(args)
	}

	if instance.file != nil {
		instance.file.DPanic(args)
	}
}

func LogPanic(args ...interface{}) {
	if instance.console != nil {
		instance.console.Error(args)
	}

	if instance.file != nil {
		instance.file.Error(args)
	}
}

func LogFatal(args ...interface{}) {
	if instance.console != nil {
		instance.console.Error(args)
	}

	if instance.file != nil {
		instance.file.Error(args)
	}
}

func LogDebugf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Debugf(template, args...)
	}

	if instance.file != nil {
		instance.file.Debugf(template, args...)
	}
}

func LogInfof(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Infof(template, args...)
	}

	if instance.file != nil {
		instance.file.Infof(template, args...)
	}
}

func LogPrintf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Infof(template, args...)
	}

	if instance.file != nil {
		instance.file.Infof(template, args...)
	}
}

func LogWarnf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Warnf(template, args...)
	}

	if instance.file != nil {
		instance.file.Warnf(template, args...)
	}
}

func LogErrf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Errorf(template, args...)
	}

	if instance.file != nil {
		instance.file.Errorf(template, args...)
	}
}

func LogDPanicf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.DPanicf(template, args...)
	}

	if instance.file != nil {
		instance.file.DPanicf(template, args...)
	}
}

//直接不走panic
func LogPanicf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Errorf(template, args...)
	}

	if instance.file != nil {
		instance.file.Errorf(template, args...)
	}
}

func LogFatalf(template string, args ...interface{}) {
	if instance.console != nil {
		instance.console.Errorf(template, args...)
	}

	if instance.file != nil {
		instance.file.Errorf(template, args...)
	}
}

func RunLog(tag string, start time.Time, timeLimit float64) {
	dis := time.Now().Sub(start).Seconds()
	if dis > timeLimit {
		Logf("%v start at %v ,cost %v ", tag,
			start.Format("2006-01-02 15:04:05"), time.Now().Sub(start))
	}
}

func SetLevel(level Level) {
	instance.atom.SetLevel(instance.ops.parseLevel(level))
}

func Close() error {
	return instance.Close()
}
