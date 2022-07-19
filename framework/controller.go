package framework

// ControllerHandler 将处理器类的函数抽象为类型
type ControllerHandler func(c *Context) error
