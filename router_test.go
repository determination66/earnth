package earnth

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAddRoute(t *testing.T) {

	testCases := []struct {
		name   string
		method string
		path   string
	}{
		{"index1", http.MethodGet, "/user"},
		{"index2", http.MethodGet, "/"},
		{"index3", http.MethodGet, "/user/home"},
		{"index4", http.MethodGet, "/order/detail"},
		{"index5", http.MethodPost, "/order/create"},
		{"index6", http.MethodPost, "/login"},
		{"index7", http.MethodGet, "/index"}, // Added index test case
	}

	//var mockHandler HandleFunc = func() {}
	//r := newRouter()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 逻辑代码
		})

	}
}
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return "path不相等", false
	}

	if len(n.children) != len(y.children) {
		return "children数目不相等", false
	}
	nHandler := reflect.ValueOf(n)
	yHandler := reflect.ValueOf(y)
	if nHandler != yHandler {
		return "Handler不相等", false
	}
	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
