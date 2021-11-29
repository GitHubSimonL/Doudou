package logger

import (
	"go.uber.org/zap"
)

//这里用于打点日志，结构化日志信息，方便查询
type KeyValue zap.Field

type EsLog struct {
	Sys       string
	User      int64
	OpType    string
	Msg       string
	ExtraInfo []KeyValue
}

func (esl *EsLog) genLog() []interface{} {
	extraLen := len(esl.ExtraInfo)

	result := make([]interface{}, extraLen+3)
	result[0] = zap.Field(GetKeyValue("Sys", esl.Sys))
	result[1] = zap.Field(GetKeyValue("User", esl.User))
	result[2] = zap.Field(GetKeyValue("OpType", esl.OpType))

	for index, data := range esl.ExtraInfo {
		result[3+index] = zap.Field(data)
	}
	return result
}

func EsLogDebug(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Debugw(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogInfo(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Infow(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogWarn(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Warnw(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogErr(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Errorw(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogDPanic(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.DPanicw(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogPanic(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Panicw(esLog.Msg, esLog.genLog()...)
	}
}

func EsLogFatal(esLog EsLog) {
	if instance.esLog != nil {
		instance.esLog.Fatalw(esLog.Msg, esLog.genLog()...)
	}
}

//生成键值对
func GetKeyValue(key string, value interface{}) KeyValue {
	return KeyValue(zap.Any(key, value))
}
