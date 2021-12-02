package _default

import (
	"Doudou/framework/itr"
	"Doudou/lib/logger"
)

type Handler struct {
	HookFuncMap map[uint32]itr.IHandle
	PoolSize    int
	TaskQueue   []chan itr.IRequest
}

func (h *Handler) DoMsgHandler(req itr.IRequest) {
	if req == nil {
		return
	}

	handle, ok := h.HookFuncMap[req.GetMsgID()]
	if !ok {
		logger.LogErrf("handle msgID:%v func not found.", req.GetMsgID())
		return
	}

	handle.PreHandle(req)
	handle.Handle(req)
	handle.AfterHandle(req)
}

func (h *Handler) RegisterHandle(msgID uint32, handle itr.IHandle) {
	panic("implement me")
}

func (h *Handler) StartWorkPool() {
	panic("implement me")
}

func (h *Handler) AddMgsToTaskPool(req itr.IRequest) {
	panic("implement me")
}

func (h *Handler) GetTaskQueueAmount() int {
	panic("implement me")
}
