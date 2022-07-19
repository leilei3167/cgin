package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	Group(string) IGroup //使得分组可嵌套
}
type Group struct {
	core   *Core  //封装core的方法,本质就是把前缀和uri组合起来
	parent *Group //如果嵌套,指向上一个Group
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{core: core, prefix: prefix, parent: nil}
}

//递归调用 获取整体的绝对uri
func (g *Group) getAbsPrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsPrefix() + g.prefix

}

func (g *Group) Group(s string) IGroup {
	return &Group{core: g.core, prefix: g.prefix + s}
}
func (g *Group) Get(s string, handler ControllerHandler) {
	s = g.getAbsPrefix() + s
	g.core.Get(s, handler)
}

func (g *Group) Post(s string, handler ControllerHandler) {
	s = g.getAbsPrefix() + s
	g.core.Post(s, handler)
}

func (g *Group) Put(s string, handler ControllerHandler) {
	s = g.getAbsPrefix() + s
	g.core.Put(s, handler)
}

func (g *Group) Delete(s string, handler ControllerHandler) {
	s = g.getAbsPrefix() + s
	g.core.Delete(s, handler)
}
