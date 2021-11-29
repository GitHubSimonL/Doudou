package logger

import (
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

type LWriter struct {
	lumberjack.Logger

	isASync bool //是否异步
	logChan chan *Buffer

	state atomic.Bool
}

const (
	LWriterStateOpen  = true  //开启
	LWriterStateClose = false //关闭
)

var _ zapcore.WriteSyncer = (*LWriter)(nil) //是否实现接口
var asyncCloseErr = make(chan error, 1)     //异步落地日志，close时，会等待最终的close结果信息

func NewLWriter(fileName string, maxSize, maxBackups, maxAge int, isCompress, isASync bool) *LWriter {
	writer := &LWriter{
		Logger: lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   isCompress,
		},
		isASync: isASync,
	}

	writer.SetState(LWriterStateOpen)

	if !isASync {
		return writer
	}

	writer.logChan = make(chan *Buffer, MaxAsyncLogSize)
	go writer.loop()

	return writer
}

func (lw *LWriter) Write(p []byte) (n int, err error) {
	if lw.GetState() != LWriterStateOpen {
		return 0, errors.New("LWriter not open.")
	}

	if !lw.isASync || len(lw.logChan) >= MaxAsyncLogSize {
		return lw.SyncWrite(p)
	}

	buff := localPool.Get()
	buff.Copy(p)

	return lw.put(buff)
}

func (lw *LWriter) SyncWrite(p []byte) (n int, err error) {
	return lw.Logger.Write(p)
}

func (lw *LWriter) put(buffer *Buffer) (n int, err error) {
	lw.logChan <- buffer

	return len(buffer.bs), nil
}

func (lw *LWriter) loop() {
	defer func() {
		asyncCloseErr <- lw.Logger.Close() //异步落地，待写完成后，关闭file
	}()

	for {
		select {
		case buffer, ok := <-lw.logChan:
			if !ok {
				return
			}

			if buffer == nil {
				continue
			}

			lw.SyncWrite(buffer.bs)
			buffer.Free()
		}
	}
}

//此方法是倍zaplogger的close调用 (log.go  instance.Close())
func (lw *LWriter) Sync() error {
	return lw.close()
}

func (lw *LWriter) close() error {
	if lw.GetState() == LWriterStateClose {
		return errors.New("LWriter already closed")
	}

	defer lw.SetState(LWriterStateClose)

	if !lw.isASync { //同步落地时直接关闭底层logger
		return lw.Logger.Close()
	}

	close(lw.logChan)

	return <-asyncCloseErr
}

//设置状态
func (lw *LWriter) SetState(state bool) {
	lw.state.Store(state)
}

//获取wirter状态，当为false时表示不在允许写日志
func (lw *LWriter) GetState() bool {
	return lw.state.Load()
}
