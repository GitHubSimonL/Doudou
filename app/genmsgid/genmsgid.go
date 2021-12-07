package main

/*
 * 导出msgid给客户端使用
 */
import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"

	"Doudou/lib/logger"
	"Doudou/protocol"
)

// msg .
type msg struct {
	ID   uint32
	Name string
}

func main() {
	path := flag.String("path", "./protocol/msgid.def", "output path")
	f, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogErrf("export msgid to %v failed!", path)
		return
	}

	defer f.Close()
	var msgList []*msg

	protocol.Processor.Range(func(id uint32, t reflect.Type) {
		msgList = append(msgList, &msg{
			ID:   id,
			Name: t.Elem().Name(),
		})
	})

	sort.Slice(msgList, func(i, j int) bool {
		return msgList[i].ID < msgList[j].ID
	})

	_ = f.Truncate(0) // 清空文件内容
	writer := bufio.NewWriter(f)
	for _, m := range msgList {
		_, _ = writer.WriteString(fmt.Sprintf("%d %s\n", m.ID, m.Name))
	}

	_ = writer.Flush()
	logger.Logf("export msgid to %s complete", *path)
}
