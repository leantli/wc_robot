package robot

import (
	"log"
)

// message_handler_chain.go 定义了消息处理链的结构与执行方法

type MsgHandlerChain struct {
	GlobalMatch IsMatchFn
	Handlers    []*Handler // 消息处理链
}

type IsMatchFn func(message *Message) bool

type HandleFn func(message *Message) error

type Handler struct {
	Name string
	// 是否匹配成功的方法
	IsMatchFn IsMatchFn
	// 处理消息的具体方法
	HandleFn HandleFn
}

// 执行处理链中的handlers
func (c *MsgHandlerChain) Handle(message *Message) {
	if !c.GlobalMatch(message) {
		return
	}
	for _, handler := range c.Handlers {
		// 匹配上了再处理，处理成了就直接返回，不继续匹配了
		if handler.IsMatchFn(message) {
			err := handler.HandleFn(message)
			if err != nil {
				log.Printf("[ERROR]处理器%s处理失败, 错误原因err:%v", handler.Name, err)
			}
			return
		}
	}
}

// 注册消息处理方法
func (c *MsgHandlerChain) RegisterHandler(name string, matchFn IsMatchFn, handleFns ...HandleFn) {
	for _, fn := range handleFns {
		h := &Handler{
			Name:      name,
			HandleFn:  fn,
			IsMatchFn: matchFn,
		}
		c.Handlers = append(c.Handlers, h)
	}
}

// 注册全局校验方法
func (c *MsgHandlerChain) RegisterGlobalCheck(matchFn IsMatchFn) {
	c.GlobalMatch = matchFn
}
