package _default

import (
	"Doudou/framework/itr"
	"Doudou/lib/logger"
	"sync"
)

type ApiMgr struct {
	ApiMap    map[uint32]itr.IHandle
	TaskQueue []chan itr.IRequest

	sync.RWMutex
}

var _ itr.IApiMgr = (*ApiMgr)(nil)

func NewApiMgr(taskPoolSize int) *ApiMgr {
	return &ApiMgr{
		ApiMap:    make(map[uint32]itr.IHandle),
		TaskQueue: make([]chan itr.IRequest, 0, taskPoolSize),
		RWMutex:   sync.RWMutex{},
	}
}

func (h *ApiMgr) RegisterHandle(msgID uint32, handle itr.IHandle) {
	if handle == nil {
		logger.LogWarnf("handle is nil. %v", msgID)
		return
	}

	h.Lock()
	defer h.Unlock()

	if _, ok := h.ApiMap[msgID]; ok {
		logger.LogWarnf("repeated register api: %v", msgID)
		return
	}

	h.ApiMap[msgID] = handle
}

func (h *ApiMgr) OneTask(taskIdx int, queue chan itr.IRequest) {
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

func (h *ApiMgr) StartWorkPool() {
	if h.GetTaskQueueAmount() < 2 {
		return
	}

	for i := 0; i < h.GetTaskQueueAmount(); i++ {
		h.TaskQueue[i] = make(chan itr.IRequest, DefaultRequestQueueLen)
		go h.OneTask(i, h.TaskQueue[i])
	}
}

func (h *ApiMgr) AddMgsToTaskPool(req itr.IRequest) {
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

func (h *ApiMgr) GetTaskQueueAmount() int {
	return cap(h.TaskQueue)
}

func (h *ApiMgr) DoMsgHandler(req itr.IRequest) {
	if req == nil {
		return
	}

	h.RLock()
	defer h.RUnlock()

	handle, ok := h.ApiMap[req.GetMsgID()]
	if !ok {
		logger.LogErrf("handle msgID:%v func not found.", req.GetMsgID())
		return
	}

	handle.PreHandle(req)
	handle.Handle(req)
	handle.AfterHandle(req)
}
