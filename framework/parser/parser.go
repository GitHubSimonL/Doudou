package parser

import (
	"Doudou/lib/bkdrhash"
	"Doudou/lib/logger"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// protobuf消息解析器
type Processor struct {
	id2type map[uint32]protoreflect.MessageType // id to type映射
	type2id map[protoreflect.MessageType]uint32 // type to id映射
}

func NewProcessor() *Processor {
	processor := new(Processor)
	processor.type2id = make(map[protoreflect.MessageType]uint32)
	processor.id2type = make(map[uint32]protoreflect.MessageType)
	return processor
}

// 注册消息类型
func (p *Processor) Register(msg proto.Message) {
	reflect := proto.MessageReflect(msg)
	msgID := bkdrhash.BKDRHash(string(reflect.Descriptor().Name()))
	// 类型必须不为空，且为指针
	if reflect.Type() == nil {
		logger.LogFatal("require protobuf message pointer")
	}
	// 不能重复注册
	if _, ok := p.type2id[reflect.Type()]; ok {
		logger.LogFatalf("message %s is already registered", reflect.Type())
	}
	// 类型和ID双向注册
	p.type2id[reflect.Type()] = msgID
	p.id2type[msgID] = reflect.Type()
}

func (p *Processor) Unmarshal(msgID uint32, data []byte) (proto.Message, error) {
	// 根据ID取类型
	typ, ok := p.id2type[msgID]
	if !ok {
		return nil, fmt.Errorf("message id %v not registered", msgID)
	}
	// 根据类型反序列化
	msg := proto.MessageV1(typ.New())
	err := proto.Unmarshal(data, msg)
	return msg, err
}

// Marshal 序列化
func (p *Processor) Marshal(msg proto.Message) (uint32, []byte, error) {
	reflect := proto.MessageReflect(msg)
	// 根据类型查找消息id
	msgId, ok := p.type2id[reflect.Type()]
	if !ok {
		err := fmt.Errorf("msg %s not registered", reflect.Type())
		return 0, nil, err
	}
	// 序列化
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, nil, err
	}
	return msgId, data, err
}

func (p *Processor) GetMsgID(msg proto.Message) uint32 {
	reflect := proto.MessageReflect(msg)

	return p.type2id[reflect.Type()]
}

func (p *Processor) GetMsgType(msgID uint32) protoreflect.MessageType {
	return p.id2type[msgID]
}

// 迭代调用
func (p *Processor) Range(f func(id uint32, t protoreflect.MessageType)) {
	for id, typ := range p.id2type {
		f(id, typ)
	}
}
