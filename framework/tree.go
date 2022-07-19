package framework

import "strings"

type Tree struct {
	root *node
}

type node struct {
	// 代表这个节点是否可以成为最终的路由规则。该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	isLast bool
	//uri中某个段的字符串
	segment string
	//这个节点包含的控制器
	handler ControllerHandler
	//这个节点下的子节点
	childs []*node
}

/*
考虑路由匹配冲突:
/user/name
/user/:id

由于存在通配符,请求其中某一个地址时两个都能匹配到,所以增加路由前 就要对存在性进行判断
判断某个uri是否存在,如果存在,返回节点指针,不存在返回nil,采用函数递归实现

首先需要将uri根据第一个 / 进行分隔,只需分为最多2个段

如果只能分为1个段 说明已经没有/ 分隔符了,此时再检查下一级节点是否有匹配这个段的节点就行

如果成功分成2个段,用第一个来检查下一级节点是否有匹配到
	-没有,说明匹配不到
	-有符合第一个的,则对所有符合的节点进行重新递归调用


*/
//判断是否有 : 通配符
func isWildSegment(segment string) bool {
	return strings.Contains(segment, ":")
}

//过滤下一层满足segment段的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	//如果segment是通配符则下一层所有子节点都满足
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))

	//遍历其所有的子节点
	for _, cnode := range n.childs {
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
	//第一部分用于匹配下一层的子节点
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

	//如果成功分为2部分,递归每个子节点
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}
