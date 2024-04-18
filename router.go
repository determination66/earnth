package earnth

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type node struct {
	path     string
	children map[string]*node

	handler HandleFunc
}
