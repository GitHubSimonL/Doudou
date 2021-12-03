package itr

type IHandle interface {
	PreHandle(request IRequest)   // 在处理conn业务之前的钩子方法
	Handle(request IRequest)      // 处理conn业务的方法
	AfterHandle(request IRequest) // 处理conn业务之后的钩子方法
}

type BaseHandle struct {
}

func (b *BaseHandle) PreHandle(request IRequest) {
}

func (b *BaseHandle) Handle(request IRequest) {
}

func (b *BaseHandle) AfterHandle(request IRequest) {
}
