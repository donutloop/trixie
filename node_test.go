package tmux

import (
	"testing"
)

func TestNodeAddType(t *testing.T) {

	pathTestCases := []struct {
		title   string
		path    string
		compare func(edge *Edge) bool
	}{
		{
			title:   "Static path",
			path:    "/echo",
			compare: func(edge *Edge) bool { return staticNode == edge.typ },
		},
		{
			title:   "Static path",
			path:    "/overview/",
			compare: func(edge *Edge) bool { return paramNode != edge.typ },
		},
		{
			title:   "Param path",
			path:    "/:string/:number/dummy",
			compare: func(edge *Edge) bool { return paramNode == edge.typ },
		},
		{
			title:   "Param path",
			path:    "/dummy/:number/dummy",
			compare: func(edge *Edge) bool { return staticNode != edge.typ },
		},
	}

	parentNode := &Node{}
	for _, pathTestCase := range pathTestCases {
		t.Run(pathTestCase.title, func(t *testing.T) {

			edge := &Edge{
				node: &Node{
					prefix: pathTestCase.path,
				},
			}

			parentNode.AddType(edge)

			if !pathTestCase.compare(edge) {
				t.Errorf("Unexpected type of node (type: %v)", edge.typ)
			}
		})
	}
}
