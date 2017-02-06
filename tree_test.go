package tmux

import (
	"net/http"
	"testing"
)

func TestRadixTreeFind(t *testing.T) {
	tree := NewRadixTree(NewNode)()

	pathTestCases := []struct {
		title   string
		method  string
		pathRaw string
		path    string
	}{
		{
			title:   "Param path",
			method:  http.MethodGet,
			pathRaw: "user/:number/comments/comment/1",
			path:    "user/1/comments/comment/1",
		},
		{
			title:   "Param path",
			method:  http.MethodGet,
			pathRaw: "user/:number/comments",
			path:    "user/1/comments",
		},
		{
			title:   "Param path",
			method:  http.MethodGet,
			pathRaw: "article/:number",
			path:    "article/4",
		},
		{
			title:   "Param path",
			method:  http.MethodGet,
			pathRaw: "article/:number/comment/:number",
			path:    "article/5/comment/6",
		},
		{
			title:   "Param path",
			method:  http.MethodGet,
			pathRaw: "/host",
			path:    "/host",
		},
		{
			title:   "Param path",
			method:  http.MethodPost,
			pathRaw: "/host",
			path:    "/host",
		},
	}

	testHandler := func() func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {}
	}

	for _, pathTestCase := range pathTestCases {

		route := NewRoute()
		route.SetPattern(pathTestCase.pathRaw)
		route.AddHandler(pathTestCase.method, http.HandlerFunc(testHandler()))

		tree.Insert(route)
	}

	for _, pathTestCase := range pathTestCases {
		t.Run(pathTestCase.title, func(t *testing.T) {

			route := tree.Find(tree.GetRoot(), Methods.lookup(pathTestCase.method), pathTestCase.path)

			if route == nil {
				t.Errorf("Route not found (Expected: %v, %v)", pathTestCase.pathRaw, pathTestCase.path)
				return
			}

			if route.GetPattern() != pathTestCase.pathRaw {
				t.Errorf("Unexpected node pattern (Expected: %v, %v, Actual: %v)", pathTestCase.pathRaw, pathTestCase.path, route.GetPattern())
			}
		})
	}
}

func TestFindNotFound(t *testing.T) {
	root := NewNode()
	tree := NewRadixTree(NewNode)()
	childNode := tree.Find(root, MethodGet, "/dummy")

	if childNode != nil {
		t.Errorf("Unexpected value (%v)", childNode)
	}
}

func BenchmarkRadixTreeFind(b *testing.B) {
	tree := NewRadixTree(NewNode)()

	paths := []string{
		"/api/user/1/comment/1",
		"/api/echo",
		"/api/user/:number/comments",
		"/api/user/donutloop",
		"/api/user/6",
		"/api/article/golang",
		"/api/article/7",
		"/api/article/97/comment/9",
		"/api/article/:number/questions/:number",
		"/api/article/:number/comment/:number/subcomment/:number",
	}

	testHandler := func() func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {}
	}

	for _, path := range paths {

		route := NewRoute()
		route.SetPattern(path)
		route.AddHandlerFunc(http.MethodGet, testHandler())
		tree.Insert(route)
	}

	for n := 0; n < b.N; n++ {
		tree.Find(tree.GetRoot(), MethodGet, "/article/3/comment/5/subcomment/5")
	}
}
