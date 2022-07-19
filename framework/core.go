package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core 相当于标准库中的Server的功能,提供路由,解析请求,响应的功能
type Core struct {
	//简单的map路由
	//按框架使用者使用路由的顺序分成四步来完善这个结构：定义路由map、注册路由、匹配路由、填充 ServeHTTP 方法。
	//router map[string]map[string]ControllerHandler //替换为前缀树路由
	router map[string]*Tree
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

//RESTful风格的路由注册,将方法名与http方法名一致,将URI全部转换为大写,注意后续匹配时也要转换为大写
//这样对外暴露就是大小写不敏感的,增加易用性

func (c *Core) Get(url string, handler ControllerHandler) {
	//统一使用大写key,避免使用时每次转换
	/*	upperUrl := strings.ToUpper(url)
		c.router["Get"][upperUrl] = handler*/

	if err := c.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Post(url string, handler ControllerHandler) {
	if err := c.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Put(url string, handler ControllerHandler) {
	if err := c.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error:", err)
	}
}

//匹配路由,没有则返回nil

func (c *Core) FindRouteByReq(req *http.Request) ControllerHandler {
	//全部统一大写
	uri := req.URL.Path
	method := req.Method
	upperMethod := strings.ToUpper(method)

	//先匹配方法
	if methodHandlers, ok := c.router[upperMethod]; ok {
		//查找第二层
		return methodHandlers.FindHandler(uri)
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
