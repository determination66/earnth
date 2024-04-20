package earnth

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
	ctx   *Context
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
func (r *router) matchRouter(method, path string) *matchInfo {
	root, ok := r.trees[method]
	if !ok {
		return nil
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}
	}
	mInfo := &matchInfo{
		n:          root,
		pathParams: map[string]string{},
	}

	//mInfo.n := root
	units := strings.Split(path[1:], "/")

	for _, unit := range units {
		next, ok := mInfo.n.children[unit]
		// can find the exact match
		if !ok {
			if mInfo.n.wildcardChild == nil && mInfo.n.colonChild == nil {
				return nil
			}
			// try to find the colon child
			if mInfo.n.colonChild != nil {
				// todo need to fix
				mInfo.pathParams[mInfo.n.colonChild.path] = unit
				mInfo.n = mInfo.n.colonChild
			}
			//tey to find the wildcard child
			if mInfo.n.wildcardChild != nil {
				mInfo.n = mInfo.n.wildcardChild
			}
		} else {
			mInfo.n = next
		}

	}
	return mInfo
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

// ParamInfo for the dynamic Params
type ParamInfo struct {
	pathParams map[string]string
	n          *node
}

// static (/user/abc) --> colon (/user/:name) --> wildcard (/user/*)
func (n *node) childOrCreate(seg string) *node {
	// special process the ":name"
	if seg[0] == ':' {
		if n.wildcardChild != nil || n.colonChild != nil {
			panic(fmt.Sprintf("router already has a wildcard or colon child[%s]", seg))
		}
		if n.colonChild == nil {
			n.colonChild = &node{
				path: seg[1:],
				//path: seg,
			}
			return n.colonChild
		}
		return n.colonChild
	}
	// special process the '*'
	if seg == "*" {
		// limit the ':'
		if n.wildcardChild != nil || n.colonChild != nil {
			panic(fmt.Sprintf("router already has a wildcard or colon child[%s]", seg))
		}
		if n.wildcardChild == nil {
			n.wildcardChild = &node{
				path: seg,
			}
			return n.wildcardChild
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

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
