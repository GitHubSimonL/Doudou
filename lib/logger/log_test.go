package logger

import (
	"runtime"
	"testing"
)

var funcList = []func(v ...interface{}){
	LogInfo,
	LogDebug,
	LogWarn,
	LogErr,
}

func TestLogger(t *testing.T) {
	InitLogger(WithFilename("/Users/Logan/Work/IGG/server/log/test.log"),
		WithLogLevel(DebugLevel))

	defer Close()
	LogDebug("debug日志", 1)
	LogInfo("info日志", 2)
	LogWarn("warn日志", 3)
	LogErr("error日志", 4)
}

func BenchmarkRandLog(t *testing.B) {
	t.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	InitLogger(WithFilename("/Users/simon/Work/IGG/server/log/test.log"),
		WithLogLevel(DebugLevel),
		WithIsCloseStdOut(true))

	// defer Close()

	t.StartTimer()
	for i := 0; i < t.N; i++ {
		funcList[i%len(funcList)]("测试日志22", "打印结果", i+1)
	}
}

// func BenchmarkRandLog2(t *testing.B) {
// 	t.ResetTimer()
// 	runtime.GOMAXPROCS(runtime.NumCPU())
//
// 	t.StartTimer()
// 	for i := 0; i < t.N; i++ {
// 		log.Print("测试日志111", "打印结果", i)
// 	}
// }
