package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Context 作为整个框架的控制器,封装req和responseWriter,对外暴露更合理的接口
type Context struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      context.Context

	handlers []ControllerHandler
	index    int //维护一个下标 控制链路移动

	//标识是否已超时
	hasTimeout bool
	//保护
	writerMux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:   r,
		response:  w,
		ctx:       r.Context(), //继承req中的Context
		writerMux: &sync.Mutex{},
		index:     -1,
	}
}
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

//#base: 实现基本函数功能,获取Req等

// WriterMux 对外暴露锁,这是一种常用的给调用方操作私有字段的方式
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

//#context:实现标准库接口,利用req中的ctx(ctx来自于server的ctx包装)轻松实现

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}
func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}
func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

//#request: 封装Req的对外接口
//一,获取查询参数

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		//返回URL中的Values
		return map[string][]string(ctx.request.URL.Query())
	}
	//否则返回空map
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			//如果存在同名的Query,返回最后一个
			inval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return inval
		}

	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}
func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

//二.获取form表单(PostForm和Form字段的区别:前者只包含Form表单,而后者同时包含Query)

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def

}
func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

//三.将req中的JSON数据快捷绑定至结构体

func (ctx *Context) BindJSON(obj any) error {
	if ctx.request != nil {
		//读出所有的数据,读取完成后body将背清空,未使的body重用,使用NopCloser将数据重新存入 便于后续处理使用
		body, err := io.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = io.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}

	} else {
		return errors.New("ctx request empty")
	}
	return nil
}

//resp: 封装resp的对外接口(便于快捷的回复对应类型的数据)

func (ctx *Context) Json(status int, obj any) error {
	if ctx.HasTimeout() {
		return nil
	}

	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.response.WriteHeader(status)

	byt, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	ctx.response.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj any, template string) error {
	return nil

}
func (ctx *Context) Text(status int, obj any, template string) error {
	return nil
}

//# 中间件控制

// Next 使中间件调用链向后移动一个,并传递ctx调用下一个处理器
// 1.用于请求的入口处,即Core的ServeHTTP方法
// 2.每个中间件的逻辑代码中都要用到
// 注意index的初始值应该为-1,这样才能保证第一次调用时值为0,即第一个控制器
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}
