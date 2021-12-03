package _default

import (
	"Doudou/framework/itr"
	"Doudou/lib/logger"
	"sync"
)

type Handler struct {
	HookFuncMap map[uint32]itr.IHandle
	PoolSize    int
	TaskQueue   []chan itr.IRequest

	sync.RWMutex
}

func (h *Handler) DoMsgHandler(req itr.IRequest) {
	if req == nil {
		return
	}

	h.RLock()
	defer h.RUnlock()

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
	if handle == nil {
		logger.LogWarnf("handle is nil. %v", msgID)
		return
	}

	h.Lock()
	defer h.Unlock()

	if _, ok := h.HookFuncMap[msgID]; ok {
		logger.LogWarnf("repeated register api: %v", msgID)
		return
	}

	h.HookFuncMap[msgID] = handle
}

func (h *Handler) OneTask(taskIdx int, queue chan itr.IRequest) {
	logger.LogDebugf("task %v start work.", taskIdx)
	defer func() {
		logger.LogDebugf("task %v work finish.")
	}()

	for {
		select {
		case req, ok := <-queue:
			if !ok {
				return
			}

			h.DoMsgHandler(req)
		}
	}
}

func (h *Handler) StartWorkPool() {
	if h.GetTaskQueueAmount() < 2 {
		return
	}

	for i := 0; i < int(h.GetTaskQueueAmount()); i++ {
		h.TaskQueue[i] = make(chan itr.IRequest, DefaultRequestQueueLen)
		go h.OneTask(i, h.TaskQueue[i])
	}
}

func (h *Handler) AddMgsToTaskPool(req itr.IRequest) {
	if req == nil {
		return
	}

	if h.GetTaskQueueAmount() < 2 {
		h.DoMsgHandler(req)
		return
	}

	idx := req.GetConnection().GetConnID() % uint32(h.GetTaskQueueAmount()) // 根据链接id做mod负载均衡
	taskQueueLen := uint32(len(h.TaskQueue))
	if idx >= taskQueueLen {
		idx = taskQueueLen - 1
	}

	h.TaskQueue[idx] <- req
}

func (h *Handler) GetTaskQueueAmount() int {
	return 2
}
