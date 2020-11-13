package gee

import (
	"fmt"
	"strings"
)

/**
使用前缀树 实现动态路由
map[string] handler 这种是静态路由


/:lang/doc
/:lang/tutorial
/:lang/intro
/about
/p/blog
/p/related
*/

type node struct {
	pattern  string  //待匹配路由,例如 /p/:lang
	part     string  //路由中的一部分, 例如 :lang
	children []*node //子节点,例如[doc,tutorial,intro]
	isWild   bool    //是否是精确匹配, part含有: 或*时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//@param pattern 路由
//@param parts 所有的子项
//@param height 树的高度. 开始为 0
func (n *node) insert(pattern string, parts []string, height int) {
	//如果到了最底层. 就会保存完整路由
	//其余节点的pattern 都为空
	if len(parts) == height {
		n.pattern = pattern //完整的路由
		return
	}
	//一次获取每层的数据
	part := parts[height]
	//查找对应的子节点
	child := n.matchChild(part)
	if child == nil { //如果子节点为空.  并添加到父节点中
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//查找路由
//@param parts所有子项
//@param height 树的层数
func (n *node) search(parts []string, height int) *node {
	//如果到了最底层  或者 子项包含*
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" { //为空表示  匹配失败
			return nil
		}
		return n //否则发挥上一层节点
	}

	//获取每层的子项
	part := parts[height]
	//查找对应的子节点
	children := n.matchChildren(part)
	//遍历所有的子节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

//查找所有完整的路由
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" { //不为空.就是完整路由
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

//精确匹配子节点
//如果是非精确. 直接返回子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//匹配素有合适的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
