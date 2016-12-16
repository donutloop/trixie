package tmux

import (
	"net/http"
	"testing"
)

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
