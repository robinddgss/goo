package goo

import (
	"fmt"
	"strings"
)

type node struct {
	path     string
	part     string
	children []*node
	isWild   bool
	handlers []HandlerFunc
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(path string, parts []string, height int, handlers []HandlerFunc) {
	if len(parts) == height {
		if n.path != "" {
			panic(fmt.Sprintf("Routing Conflict %s and %s\n", n.path, path))
		}
		n.path = path
		n.handlers = handlers
		return
	}
	part := parts[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == '*' || part[0] == ':',
		}
		n.children = append(n.children, child)
	}
	child.insert(path, parts, height+1, handlers)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
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
