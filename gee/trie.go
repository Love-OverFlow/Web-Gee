package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang  只有根节点不为空
	part     string  // 路由中的一部分，例如 :lang  保存路径, 这里和经典前缀树稍有不同, 路径是保存在节点中的
	children []*node // 子节点，例如 [doc, tutorial, intro]  因为路径上不确定有多少种字符, 所以是切片
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为 true
}

// 第一个匹配成功的节点，用于插入 insert
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找 search
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 类似广度优先遍历, 找到任意一个符号条件的路径就返回
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
