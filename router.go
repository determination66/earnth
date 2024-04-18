package earnth

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" {
		panic(fmt.Sprintf("path can't be empty"))
	}

	// find the tree,if the tree is nil,then create
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path == "/" {
		// 根节点重复注册
		if root.handler != nil {
			panic("web: router already has a handler[/]")
		}
		root.handler = handleFunc
		return
	}

	segs := strings.Split(path[1:], "/")

	for _, seg := range segs {
		if seg == "" {
			panic("web: 不能有连续的 /")
		}
		child := root.childOrCreate(seg)
		root = child
	}
	//if root.handler != nil {
	//	panic(fmt.Sprintf("web: 路由冲突，重复注册[%s]"))
	//}
	//root.handler = handler
}

func (r *router) match(method, path string) *node {
	root, ok := r.trees[method]
	if !ok {
		return nil
	}
	if path == "/" {
		return root
	}
	current := root
	units := strings.Split(path[1:], "/")
	for _, unit := range units {
		next, ok := current.children[unit]
		if !ok {
			return nil
		}
		current = next
	}
	return current
}

func (r *router) isEqual(y *router) bool {
	if len(r.trees) != len(y.trees) {
		return false
	}
	for method, root := range r.trees {
		dst, ok := y.trees[method]
		if !ok {
			return false
		}
		if !root.isEqual(dst) {
			return false
		}
	}
	return true
}

type node struct {
	path     string
	children map[string]*node

	handler HandleFunc
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}

	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func (n *node) isEqual(y *node) bool {
	if n == nil && y == nil {
		return true
	}
	if n == nil || y == nil {
		return false
	}
	if n.path != y.path {
		return false
	}
	if len(n.children) != len(y.children) {
		return false
	}

	for key, child := range n.children {
		dst, ok := y.children[key]
		if !ok {
			return false
		}
		if !child.isEqual(dst) {
			return false
		}
	}
	return true
}
