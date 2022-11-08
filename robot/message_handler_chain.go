package robot

import "log"

// message_handler_chain.go 定义了消息处理链的结构与执行方法

type MsgHandlerChain struct {
	Handlers []*Handler // 消息处理链
}

type HandleFn func(message *Message) error

type Handler struct {
	Name string
	// 处理消息的具体方法
	HandleFn HandleFn
}

// 执行处理链中的handlers
func (c *MsgHandlerChain) Handle(message *Message) {
	for _, handler := range c.Handlers {
		err := handler.HandleFn(message)
		if err != nil {
			log.Printf("[ERROR]处理器%s处理失败, 错误原因err:%v", handler.Name, err)
		}
	}
}

// 注册消息处理方法
func (c *MsgHandlerChain) RegisterHandler(name string, handleFns ...HandleFn) {
	for _, fn := range handleFns {
		h := &Handler{
			Name:     name,
			HandleFn: fn,
		}
		c.Handlers = append(c.Handlers, h)
	}
}
