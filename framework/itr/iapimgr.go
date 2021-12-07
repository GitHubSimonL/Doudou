package itr

type HandleFunc func(request IRequest)

type IApiMgr interface {
	RegisterHandle(msgID uint32, fn HandleFunc) // 新增处理函数
	StartWorkPool()                             // 开启工作线程，当为多线程时需自行实现根据request的id去做负载均衡
	AddMgsToTaskPool(req IRequest)              // 请求发布到任务池
}
