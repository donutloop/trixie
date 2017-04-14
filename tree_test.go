package trixie

import (
	"testing"
)

var treeRouteTestCases = []routeTestCase{
	{
		rawPath:      "/home/user/comment/sub",
		path:         "/home/user/comment/sub",
		countOfParam: 4,
	},
	{
		rawPath:      "/home/user/comment",
		path:         "/home/user/comment",
		countOfParam: 3,
	},

	{
		rawPath:      "/home/user/article/comment",
		path:         "/home/user/article/comment",
		countOfParam: 4,
	},
	{
		rawPath:      "/home/user/article/comment/:string",
		path:         "/home/user/article/comment/test",
		countOfParam: 5,
	},
	{
		rawPath:      "/:string/:string/:string/:string/:string",
		path:         "/dummy/dummy/dummy/dummy/dummy",
		countOfParam: 5,
	},
	{
		rawPath:      "/:number/:number/:number/:number/:number",
		path:         "/1/1/1/1/1",
		countOfParam: 5,
	},
	{
		rawPath:      "/#([0-9]{3,})/:number/:number/:number/:number",
		path:         "/140/1/1/1/1",
		countOfParam: 5,
	},
}

func TestTree_Insert_and_find(t *testing.T) {

	tree := NewTree(NewNode)()

	for _, testCase := range treeRouteTestCases {
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

	for _, testCase := range treeRouteTestCases {
		route, params := tree.Find(tree.GetRoot(), MethodGet, testCase.path)

		if route == nil {
			t.Errorf("Unexpected nil route (Expected: %s)", testCase.path)
			return
		}

		if testCase.countOfParam != len(params) {
			t.Errorf("Count of parameters is bad (Actual: %d, Expected: %d)", testCase.countOfParam, len(params))
			return
		}

		if route.GetPattern() != testCase.rawPath {
			t.Errorf("Unexpected route (Expected: %s, Actual: %s)", testCase.rawPath, route.GetPattern())
			return
		}
	}
}

func BenchmarkTree_Insert(b *testing.B) {
	tree := NewTree(NewNode)()
	for n := 0; n < b.N; n++ {
		for _, testCase := range treeRouteTestCases {
			route := NewRoute()
			route.SetPattern(testCase.rawPath)
			route = tree.Insert(route)
		}
	}
}

func BenchmarkTestTree_find(b *testing.B) {

	tree := NewTree(NewNode)()

	for _, testCase := range treeRouteTestCases {
		route := NewRoute()
		route.SetPattern(testCase.rawPath)
		route = tree.Insert(route)
	}

	for _, testCase := range treeRouteTestCases {
		b.Run("", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				tree.Find(tree.GetRoot(), MethodGet, testCase.path)
			}
		})
	}
}
