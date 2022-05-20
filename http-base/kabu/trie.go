package kabu

import "strings"

//定义一个子节点
type node struct {
	pattern  string  //待匹配路由
	part     string  //部分路由
	children []*node //子节点
	isWild   bool    //是否精准匹配，即是否有：或*
}

//第一个成功匹配的节点 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child) //判断是否精准匹配
	}
	child.insert(pattern, parts, height+1) //递归进行插入直到匹配结束
}

//查找
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //查找到最后一个匹配项时结束递归
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part) //找到符合的所有孩子

	for _, child := range children {
		result := child.search(parts, height+1) //递归进行查找
		if result != nil {
			return result
		}
	}
	return nil
}
