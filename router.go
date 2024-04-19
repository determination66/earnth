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
		if root.handler != nil {
			panic("web: router already has a handler[/]")
		}
		root.handler = handleFunc
		return
	}

	secs := strings.Split(path[1:], "/")

	for _, seg := range secs {
		if seg == "" {
			panic(fmt.Sprintf("Routes cannot have consecutive '/'"))
		}
		// This is important
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic(fmt.Sprintf("Duplicate router: [method:%s path:%s]", method, path))
	}
	root.handler = handleFunc
	fmt.Println("add", method, path)
}

// static routers first, then dynamic routers
// static (/user/abc) --> colon (/user/:name) --> wildcard (/user/*)
func (r *router) matchRouter(method, path string) *node {
	root, ok := r.trees[method]
	if !ok {
		return nil
	}
	if path == "/" {
		return root
	}
	cur := root
	units := strings.Split(path[1:], "/")

	for _, unit := range units {
		next, ok := cur.children[unit]
		// can't find the exact match
		if !ok {
			// try to find the colon child
			if cur.handler != nil {
				// todo need to fix

				cur = cur.colonChild

			} else {
				//tey to find the wildcard child
				if cur.wildcardChild != nil {
					cur = cur.wildcardChild
				}
			}
			return nil
		}
		cur = next
	}
	return cur
}

type node struct {
	path string
	// static router match
	children map[string]*node

	//wildcard router node ,to parse '*'
	wildcardChild *node

	//colon router node ,to parse ":name"
	colonChild *node
	//paramsChild *ParamNode

	handler HandleFunc
}

// ParamNode router node ,to parse ":name"
//type ParamNode struct {
//	name string //parameter name
//	n    *node
//}

// ParamInfo for the dynamic Params
type ParamInfo struct {
	pathParams map[string]string
	n          *node
}

// static (/user/abc) --> colon (/user/:name) --> wildcard (/user/*)
func (n *node) childOrCreate(seg string) *node {
	// special process the ":name"
	if seg[0] == ':' {
		if n.colonChild == nil {
			n.colonChild = &node{
				path: seg[1:],
			}
			return n.colonChild
		}
		return n.colonChild
	}
	// special process the '*'
	if seg == "*" {
		n.wildcardChild = &node{
			path: seg,
		}
		return n.wildcardChild
	}
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
