// 在非Windows平台编译
// +build !windows

package logger

import "os"
import "strings"
import "syscall"

func redirectStdErrLog() error {
	panicFile := strings.Replace(instance.fileName, ".log", ".panic", -1)
	fd, err := os.OpenFile(panicFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	err = syscall.Dup2(int(fd.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return err
	}

	LogDebug("=================redirect std err log success======================")

	return nil
}
