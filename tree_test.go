package tmux

import (
	"net/http"
	"testing"
)

func TestRadixTreeFind(t *testing.T) {
	tree := NewRadixTree(NewNode, NewRoute)()

	pathTestCases := []struct {
		title   string
		method  Method
		pathRaw string
		path    string
	}{
		{
			title:   "Param path",
			method:  MethodGet,
			pathRaw: "user/:number/comments/comment/1",
			path:    "user/1/comments/comment/1",
		},
		{
			title:   "Param path",
			method:  MethodGet,
			pathRaw: "user/:number/comments",
			path:    "user/1/comments",
		},
		{
			title:   "Param path",
			method:  MethodGet,
			pathRaw: "article/:number",
			path:    "article/4",
		},
		{
			title:   "Param path",
			method:  MethodGet,
			pathRaw: "article/:number/comment/:number",
			path:    "article/5/comment/6",
		},
	}

	testHandler := func() func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {}
	}

	for _, pathTestCase := range pathTestCases {
		tree.Insert(pathTestCase.method, pathTestCase.pathRaw, http.HandlerFunc(testHandler()))
	}

	for _, pathTestCase := range pathTestCases {
		t.Run(pathTestCase.title, func(t *testing.T) {

			route := tree.Find(tree.GetRoot(), pathTestCase.method, pathTestCase.path)

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
	tree := NewRadixTree(NewNode, NewRoute)()
	childNode := tree.Find(root, MethodGet, "/dummy")

	if childNode != nil {
		t.Errorf("Unexpected value (%v)", childNode)
	}
}

func BenchmarkRadixTreeFind(b *testing.B) {
	tree := NewRadixTree(NewNode, NewRoute)()

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
		tree.Insert(MethodGet, path, http.HandlerFunc(testHandler()))
	}

	for n := 0; n < b.N; n++ {
		tree.Find(tree.GetRoot(), MethodGet, "/article/3/comment/5/subcomment/5")
	}
}
