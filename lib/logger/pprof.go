package logger

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

const TTL = 5 * time.Minute

type Task struct {
	sync.Once
	LastUpdate time.Time
	Timer      *time.Timer
}

//Usages: go tool pprof xxx.prof

var pt *Task

func init() {
	pt = new(Task)
	pt.LastUpdate = time.Now()
}

func (t *Task) loop() {
	defer func() {
		if err := recover(); err != nil {
			LogErrf("pprof catch err:%v", err)
		}

		if t.Timer == nil {
			return
		}
		t.Timer.Stop()
	}()

	for {
		select {
		case <-t.Timer.C:
			if time.Now().Sub(t.LastUpdate) >= TTL {
				t.DumpGoroutineInfo()
				t.DumpCPUInfo()
				t.DumpHeapInfo()
				return
			}
			t.Timer.Reset(TTL)
		}
	}
}

func (t *Task) GetFileName(category string) string {
	fileName := instance.fileName
	fileName = strings.Replace(fileName, ".log", "", 1)

	return fmt.Sprintf("%s-%s-%d.prof", fileName, category, time.Now().Unix())
}

func (t *Task) DumpGoroutineInfo() {
	p := pprof.Lookup("goroutine")
	fileName := t.GetFileName("goroutine")

	if f, err := os.Create(fileName); err != nil {
		LogDebugf("lookup goroutine profile failed: %v", err)
	} else {
		LogDebugf("lookup goroutine profile")
		p.WriteTo(f, 0)
	}
}

func (t *Task) DumpHeapInfo() {
	p := pprof.Lookup("heap")
	fileName := t.GetFileName("heap")

	if f, err := os.Create(fileName); err != nil {
		LogDebugf("lookup heap profile failed: %v", err)
	} else {
		LogDebugf("lookup heap profile")
		p.WriteTo(f, 0)
	}
}

func (t *Task) DumpCPUInfo() {
	go func() {
		fileName := t.GetFileName("cpu")
		f, err := os.Create(fileName)
		if err != nil {
			return
		}

		timeDuration := 1 * time.Minute
		tmpTimer := time.NewTimer(timeDuration)

		defer func() {
			pprof.StopCPUProfile()
			f.Close()
			tmpTimer.Stop()
		}()

		pprof.StartCPUProfile(f)

		for {
			select {
			case <-tmpTimer.C:
				return
			}
		}
	}()
}

func StartPprofWork() {
	pt.Once.Do(func() {
		pt.Timer = time.NewTimer(TTL)
		go pt.loop()
	})
}

func PprofKeepAlive() {
	if pt == nil || pt.Timer == nil {
		return
	}
	pt.LastUpdate = time.Now()
}
