package earnth

import (
	"net/http"
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

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.AddRoute(tc.method, tc.path, mockHandler)
		})
	}

}

func TestRouter_match(t *testing.T) {
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	// 添加一些路由
	r.AddRoute(http.MethodGet, "/user", mockHandler)
	r.AddRoute(http.MethodGet, "/", mockHandler)
	r.AddRoute(http.MethodGet, "/user/home", mockHandler)
	r.AddRoute(http.MethodGet, "/order/detail", mockHandler)
	r.AddRoute(http.MethodPost, "/order/create", mockHandler)
	r.AddRoute(http.MethodPost, "/login", mockHandler)
	r.AddRoute(http.MethodGet, "/index", mockHandler)
	r.AddRoute(http.MethodGet, "/index/*", mockHandler)
	r.AddRoute(http.MethodGet, "/index/*/add", mockHandler)
	r.AddRoute(http.MethodGet, "/index/:name/add", mockHandler)
	r.AddRoute(http.MethodGet, "/index/:detail/del", mockHandler)

	testCases := []struct {
		name     string
		method   string
		path     string
		expected *node // 期望匹配的节点
	}{
		{"existing route", http.MethodGet, "/user", r.trees[http.MethodGet].children["user"]},
		{"root route", http.MethodGet, "/", r.trees[http.MethodGet]},
		{"nested route", http.MethodGet, "/user/home", r.trees[http.MethodGet].children["user"].children["home"]},
		{"non-existing route", http.MethodGet, "/notfound", nil},
		{"existing route with params", http.MethodGet, "/index", r.trees[http.MethodGet].children["index"]},
		{"non-existing method", http.MethodPut, "/user", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := r.matchRouter(tc.method, tc.path)
			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestRouter_IsEqual(t *testing.T) {
	r1 := newRouter()
	r2 := newRouter()

	var mockHandler HandleFunc = func(ctx *Context) {}

	// 添加一些路由到 r1
	r1.AddRoute("GET", "/user", mockHandler)
	r1.AddRoute("POST", "/order", mockHandler)

	// 添加相同的路由到 r2
	r2.AddRoute("GET", "/user", mockHandler)
	r2.AddRoute("POST", "/order", mockHandler)

	//添加不同的路由到 r2
	r2.AddRoute("PUT", "/product", mockHandler)

	// 检查相等的情况
	if !r1.isEqual(r2) {
		t.Errorf("Expected r1 and r2 to be equal")
	}

	// 删除 r2 中的一个路由
	delete(r2.trees, "PUT")

	// 检查不相等的情况
	if r1.isEqual(r2) {
		t.Errorf("Expected r1 and r2 to be different")
	}
}
