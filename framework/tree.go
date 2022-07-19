package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: newNode()}
}

type node struct {
	// 代表这个节点是否可以成为最终的路由规则。该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	isLast bool
	//uri中某个段的字符串（根节点设为空）
	segment string
	//这个节点包含的控制器
	handler ControllerHandler
	//这个节点下的子节点（多叉结构）
	childs []*node
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

/*
考虑路由匹配冲突:
/user/name
/user/:id
以上两个路径实际就冲突了,由于存在通配符,请求其中某一个地址时两个都能匹配到,所以增加路由前 就要对存在性进行判断

我们可以用 matchNode 方法，寻找某个路由在 trie 树中匹配的节点，如果有匹配节点，返回节点指针，
否则返回 nil。matchNode 方法的参数是一个 URI，返回值是指向 node 的指针，它的实现思路是使用函数递归,思路:

1.首先需要将uri根据第一个 / 进行分隔,只需分为最多2个段

2.如果只能分为1个段 说明已经没有/ 分隔符了,此时再检查下一级节点是否有匹配这个段的节点就行

3.如果成功分成2个段,用第一个来检查下一级节点是否有匹配到
	-没有,说明匹配不到
	-有符合第一个的,则对所有符合的节点进行重新递归调用

整个流程会频繁使用过滤下一层满足segment规则的子节点,所以将其提取,其逻辑就是遍历下一层子节点,判断segment与传入的是否匹配

*/

//# 首先定义树的结构,以及判断是否匹配的逻辑,后续增加等操作都会基于判断匹配
//判断是否有 : 通配符
func isWildSegment(segment string) bool {
	return strings.Contains(segment, ":")
}

//过滤某个node下一层满足segment段的子节点(包含通配符子节点)
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	//如果segment是通配符则下一层所有子节点都满足
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))

	//非通配符,遍历其所有的子节点,对比segment
	for _, cnode := range n.childs {
		//子节点是通配符,则他肯定匹配
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment { //子节点的segment匹配到
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

//判断是否已经在树中存在
func (n *node) matchNode(uri string) *node {

	segments := strings.SplitN(uri, "/", 2)
	//第一部分用于匹配下一层的子节点(将前缀用于匹配)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	//匹配符合的下一层节点
	// 如果当前子节点没有一个符合，那么说明这个uri一定是之前不存在, 直接返回nil
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	//如果能Split成1个segment,则说明是最后一个
	if len(segments) == 1 {
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		return nil
	}

	//如果成功分为2部分,递归每个子节点(使用第一个前缀后面的部分)
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

/*
增加路由的逻辑:
1.首先确认路由是否冲突,先检查要增加的规则已经在树中有可以匹配的节点了,会发生冲突,返回错误
2.增加路由的每个段时,先去树的每一层匹配查找,如果有符合这个段的节点,就不创建该节点,继续匹配待增加路由的下个段,否则需要新建一个节点
来代表这个段,这里使用到filterChildNodes

/book/list
/book/:id (冲突)
/book/:id/name
/book/:student/age
/:user/name
/:user/name/:age(冲突)
*/

func (t *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := t.root
	//是否冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exist:" + uri)
	}

	segments := strings.Split(uri, "/")
	//对每个段,依次进行匹配
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1 //便于识别每个segment是否是最后一段

		var objNode *node //标记是否有合适的子节点

		childNodes := n.filterChildNodes(segment)
		//如果有匹配的子节点,则选择这个子节点,没有则直接创建新的
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		//没有则创建一个节点容纳segment
		if objNode == nil {
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}
	return nil
}

/*
查找路由的逻辑:
*/

func (t *Tree) FindHandler(uri string) ControllerHandler {
	matchNode := t.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}

/*
最后 将增加 和查找路由 功能添加到框架中,将Core结构中map路由表替换为树
*/
