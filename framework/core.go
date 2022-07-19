package framework

import (
	"net/http"
	"strings"
)

// Core 相当于标准库中的Server的功能,提供路由,解析请求,响应的功能
type Core struct {
	//简单的map路由
	//按框架使用者使用路由的顺序分成四步来完善这个结构：定义路由map、注册路由、匹配路由、填充 ServeHTTP 方法。
	router map[string]map[string]ControllerHandler
}

func NewCore() *Core {
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}

	//写入一级map
	router := map[string]map[string]ControllerHandler{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter
	return &Core{router: router}
}

//RESTful风格的路由注册,将方法名与http方法名一致,将URI全部转换为大写,注意后续匹配时也要转换为大写
//这样对外暴露就是大小写不敏感的,增加易用性

func (c *Core) GET(url string, handler ControllerHandler) {
	//统一使用大写key,避免使用时每次转换
	upperUrl := strings.ToUpper(url)
	c.router["GET"][upperUrl] = handler
}

func (c *Core) POST(url string, handler ControllerHandler) {
	//统一使用大写key,避免使用时每次转换
	upperUrl := strings.ToUpper(url)
	c.router["POST"][upperUrl] = handler
}

func (c *Core) PUT(url string, handler ControllerHandler) {
	//统一使用大写key,避免使用时每次转换
	upperUrl := strings.ToUpper(url)
	c.router["PUT"][upperUrl] = handler
}

func (c *Core) DELETE(url string, handler ControllerHandler) {
	//统一使用大写key,避免使用时每次转换
	upperUrl := strings.ToUpper(url)
	c.router["DELETE"][upperUrl] = handler
}

//匹配路由,没有则返回nil

func (c *Core) FindRouteByReq(req *http.Request) ControllerHandler {
	//全部统一大写
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)
	upperUri := strings.ToUpper(uri)

	//先匹配方法
	if methodHandlers, ok := c.router[upperMethod]; ok {
		//查找第二层
		if handler, ok := methodHandlers[upperUri]; ok {
			return handler
		}
	}
	return nil
}

/*
当然可以新建Group结构体来承载一样的方法,但是考虑一下后续如果Group实现要修改该怎么办
更好的办法是使用接口来替代结构体的定义! 如果返回的是iGroup接口,后续有改动不需要修改
具体实例的定义,只需修改实现的方法即可

type Group struct {
	core
	prefix string
}

func (c *Core) Group() *Group {

}

选择使用接口,是考虑到未来的拓展性,框架设计层面
如果你觉得这个模块是完整的，而且后续希望有扩展的可能性，那么就应该尽量使用接口来替代实现

*/

type IGroup interface {
	GET(string, ControllerHandler)
	POST(string, ControllerHandler)
	PUT(string, ControllerHandler)
	DELETE(string, ControllerHandler)
}
type Group struct {
	core   *Core //封装core的方法,本质就是把前缀和uri组合起来
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{core: core, prefix: prefix}
}

func (g *Group) GET(s string, handler ControllerHandler) {
	s = g.prefix + s
	g.core.GET(s, handler)
}

func (g *Group) POST(s string, handler ControllerHandler) {
	s = g.prefix + s
	g.core.POST(s, handler)
}

func (g *Group) PUT(s string, handler ControllerHandler) {
	s = g.prefix + s
	g.core.PUT(s, handler)
}

func (g *Group) DELETE(s string, handler ControllerHandler) {
	s = g.prefix + s
	g.core.DELETE(s, handler)
}
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c Core) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//封装自定义的context
	ctx := NewContext(request, writer)

	//查找路由
	router := c.FindRouteByReq(request)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}
	//调用查找到的处理器处理

	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}

}
