package tmux

import "testing"

func TestReplaceEdge(t *testing.T) {

	parentNode := &Node{
		edges: [edgeTypes]Edges{},
	}

	edges := []*Edge{
		&Edge{
			label: "a"[0],
			typ:   paramNode,
			node: &Node{
				prefix: "artcile/:number",
			},
		},
		&Edge{
			label: ":"[0],
			typ:   paramNode,
			node: &Node{
				prefix: ":number",
			},
		},
		&Edge{
			label: "u"[0],
			typ:   staticNode,
			node: &Node{
				prefix: "user",
			},
		},
	}

	for _, edge := range edges {
		parentNode.AddEdge(edge)
	}

	newEdge := &Edge{
		label: "a"[0],
		typ:   staticNode,
		node: &Node{
			prefix: "article/",
		},
	}

	if err := parentNode.ReplaceEdge(newEdge); err != nil {
		t.Errorf("Unexpected error while replaceing (%s)", err.Error())
	}
}

func BenchmarkGetEdge(b *testing.B) {

	parentNode := &Node{
		edges: [edgeTypes]Edges{},
	}

	edges := []*Edge{
		&Edge{
			label: "a"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "b"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "c"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "d"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "e"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "f"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "g"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "h"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "i"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "j"[0],
			typ:   staticNode,
		},
		&Edge{
			label: "k"[0],
			typ:   staticNode,
		},
	}

	for _, edge := range edges {
		parentNode.edges[edge.typ] = append(parentNode.edges[edge.typ], edge)
	}
	parentNode.edges[staticNode].Sort()

	for _, v := range []string{"a", "f", "j"} {
		b.Run(v, func(b *testing.B) {

			for i := 0; i <= b.N; i++ {
				parentNode.GetEdge(v[0])
			}
		})
	}
}

func TestReplaceEdgeFail(t *testing.T) {

	parentNode := &Node{
		edges: [edgeTypes]Edges{},
	}

	newEdge := &Edge{
		label: "a"[0],
		typ:   staticNode,
		node: &Node{
			prefix: "article/",
		},
	}

	if err := parentNode.ReplaceEdge(newEdge); err == nil {
		t.Errorf("Unexpected success while replaceing")
	}
}

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
