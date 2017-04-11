package tmux

import (
	"testing"
)

func TestTree_Insert(t *testing.T) {

	testCases := []struct {
		rawPath string
		path    string
	}{
		{
			rawPath: "/home/user/comment/sub",
			path:    "/home/user/comment/sub",
		},
		{
			rawPath: "/home/user/comment",
			path:    "/home/user/comment",
		},

		{
			rawPath: "/home/user/article/comment",
			path:    "/home/user/article/comment",
		},
		{
			rawPath: "/home/user/article/comment/:string",
			path:    "/home/user/article/comment/test",
		},
		{
			rawPath: "/:string/:string/:string/:string/:string",
			path:    "/dummy/dummy/dummy/dummy/dummy",
		},
		{
			rawPath: "/:number/:number/:number/:number/:number",
			path:    "/1/1/1/1/1",
		},
		{
			rawPath: "/#([0-9]{3,})/:number/:number/:number/:number",
			path:    "/140/1/1/1/1",
		},
	}

	tree := NewTree(NewNode)()

	for _, testCase := range testCases {
		route := NewRoute()
		route.SetPattern(testCase.rawPath)
		route = tree.Insert(route)

		if route == nil {
			t.Errorf("Unexpected nil route (Expected: %s)", testCase.rawPath)
			return
		}

		if route.GetPattern() != testCase.rawPath {
			t.Errorf("Unexpected route (Expected: %s, Actual: %s)", testCase.rawPath, route.GetPattern())
		}
	}

	for _, testCase := range testCases {
		route, _ := tree.Find(tree.GetRoot(), MethodGet, testCase.path)

		if route == nil {
			t.Errorf("Unexpected nil route (Expected: %s)", testCase.path)
			return
		}

		if route.GetPattern() != testCase.rawPath {
			t.Errorf("Unexpected route (Expected: %s, Actual: %s)", testCase.rawPath, route.GetPattern())
		}
	}

}
