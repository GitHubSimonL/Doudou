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

// MsgList 用于排序
type MsgList []*msg

// Len 长度
func (ml MsgList) Len() int {
	return len(ml)
}

// Swap 交换i, j
func (ml MsgList) Swap(i, j int) {
	ml[i], ml[j] = ml[j], ml[i]
}

// Less elem(i) < elem(j)
func (ml MsgList) Less(i, j int) bool {
	return ml[i].ID < ml[j].ID
}

func main() {
	path := flag.String("path", "./src/protocol/msgid.def", "output path")
	f, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogErrf("export msgid to %v failed!", path)
		return
	}
	defer f.Close()
	var msglist MsgList
	protocol.Processor.Range(func(id uint32, t reflect.Type) {
		// logger.Infof("export: %s", t.Elem().MailTitle())
		msglist = append(msglist, &msg{
			ID:   id,
			Name: t.Elem().Name(),
		})
	})
	sort.Sort(msglist)
	_ = f.Truncate(0) // 清空文件内容
	writer := bufio.NewWriter(f)
	for _, m := range msglist {
		_, _ = writer.WriteString(fmt.Sprintf("%d %s\n", m.ID, m.Name))
	}
	_ = writer.Flush()
	logger.Logf("export msgid to %s complete", *path)
}
