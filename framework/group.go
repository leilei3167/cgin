package framework

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	Group(string) IGroup //使得分组可嵌套

	Use(middlewares ...ControllerHandler)
}
type Group struct {
	core   *Core  //封装core的方法,本质就是把前缀和uri组合起来
	parent *Group //如果嵌套,指向上一个Group
	prefix string

	middlewares []ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

//递归调用 获取整体的绝对uri
func (g *Group) getAbsPrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsPrefix() + g.prefix
}

func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}
func (g *Group) Get(s string, handlers ...ControllerHandler) {
	s = g.getAbsPrefix() + s
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(s, allHandlers...)
}

func (g *Group) Post(s string, handlers ...ControllerHandler) {
	s = g.getAbsPrefix() + s
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Post(s, allHandlers...)
}

func (g *Group) Put(s string, handlers ...ControllerHandler) {
	s = g.getAbsPrefix() + s
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Put(s, allHandlers...)
}

func (g *Group) Delete(s string, handlers ...ControllerHandler) {
	s = g.getAbsPrefix() + s
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Delete(s, allHandlers...)
}

// #对某个组添加中间件

// 注册中间件

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}
