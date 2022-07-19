package framework

import (
	"log"
	"net/http"
)

// Core 相当于标准库中的Server的功能,提供路由,解析请求,响应的功能
type Core struct {
	//简单的map路由
	router map[string]ControllerHandler
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}
func (c Core) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, writer)

	//写死,为了测试
	router := c.router["foo"]
	if router == nil {
		return
	}
	log.Println("core.router")
	router(ctx)
}
