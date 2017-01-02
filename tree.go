package tmux

import (
	"regexp"
	"strings"
)

// RouteTreeInterface if like you to implement your own tree version, feel free to do it
type RouteTreeInterface interface {
	UseNode(func() NodeInterface)
	Insert(RouteInterface) RouteInterface
	Find(NodeInterface, Method, string) RouteInterface
	GetRoot() NodeInterface
}

// RadixTree implements RouteTreeInterface. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and ordered iteration.
// based on go-radix ideas (github.com/armon/go-radix)
type RadixTree struct {
	root            NodeInterface
	nodeConstructor func() NodeInterface
}

// NewRadixTree returns an empty Radix Tree
func NewRadixTree(nodeConstructor func() NodeInterface) func() RouteTreeInterface {
	return func() RouteTreeInterface {
		tree := &RadixTree{}
		tree.UseNode(nodeConstructor)
		tree.root = tree.nodeConstructor()
		return tree
	}
}

// UseNode that you can use diffrent node versions
// See NodeInterface for more details (node.go)
func (t *RadixTree) UseNode(constructer func() NodeInterface) {
	t.nodeConstructor = constructer
}

func (t *RadixTree) GetRoot() NodeInterface {
	return t.root
}

// Insert is used to add a new entry or update
// an existing entry.
func (t *RadixTree) Insert(newRoute RouteInterface) RouteInterface {
	var parent NodeInterface
	currentNode := t.root
	search := newRoute.GetPattern()

	for {
		// Handle key exhaution
		if len(search) == 0 {
			if currentNode.IsLeaf() {
				currentNode.SetLeaf(mergeRoutes(currentNode.GetLeaf(), newRoute))
				return currentNode.GetLeaf()
			}

			currentNode.SetLeaf(newRoute)
			return currentNode.GetLeaf()
		}

		// Look for the edge
		currentEdge := currentNode.GetEdge(search[0])

		parent = currentNode

		// No edge, create one
		if currentEdge == nil {
			newNode := t.nodeConstructor().SetPrefixPath(search).SetLeaf(newRoute)

			parent.AddEdge(&Edge{
				label: search[0],
				node:  newNode,
			})
			return newRoute
		}

		currentNode = currentEdge.node

		// Determine longest prefix of the search key on currentNode
		commonPrefix := longestPrefix(search, currentNode.GetPrefixPath())

		// Check if they share the same prefix when yes overwrite current search and continue to next iteration
		if commonPrefix == len(currentNode.GetPrefixPath()) {
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		childNode := t.nodeConstructor().SetPrefixPath(search[:commonPrefix])

		err := parent.ReplaceEdge(&Edge{
			label: search[0],
			node:  childNode,
		})

		if err != nil {
			panic(err.Error())
		}

		// Restore the existing node
		childNode.AddEdge(&Edge{
			label: currentNode.GetPrefixPath()[commonPrefix],
			node:  currentNode.SetPrefixPath(currentNode.GetPrefixPath()[commonPrefix:]),
		})

		// If the new key is a subset, add to to this node
		search = search[commonPrefix:]

		if len(search) == 0 {
			childNode.SetLeaf(newRoute)
			return newRoute
		}

		// Create a new edge for the node
		newEdgeNode := t.nodeConstructor().SetPrefixPath(search).SetLeaf(newRoute)

		childNode.AddEdge(&Edge{
			label: search[0],
			node:  newEdgeNode,
		})

		return newRoute
	}
}

// Find is used to lookup a specific key, returning
// the value and if it was found
func (t *RadixTree) Find(root NodeInterface, method Method, path string) RouteInterface {

	for typ, edges := range root.GetEdges() {

		if len(edges) == 0 {
			continue
		}

		node, pathSegment, ok := t.findEdge(egdeType(typ), edges, path)

		if !ok {
			continue
		}

		if len(pathSegment) == 0 && node.IsLeaf() {
			return node.GetLeaf()
		}

		// recursively find the next node.
		route := t.Find(node, method, pathSegment)
		if route != nil {
			return route
		}
	}

	return nil
}

func (t *RadixTree) findEdge(typ egdeType, edges Edges, path string) (NodeInterface, string, bool) {

	var matcher func(edge *Edge) (string, bool)

	switch typ {
	case staticNode:
		matcher = func(edge *Edge) (string, bool) {

			if edge.label == byte(path[0]) && strings.HasPrefix(path, edge.node.GetPrefixPath()) {
				return path[len(edge.node.GetPrefixPath()):], true
			}

			return "", false
		}
	case paramNode:
		fallthrough
	case regexNode:
		matcher = func(edge *Edge) (string, bool) {
			reg := regexp.MustCompile(edge.node.GetPrefixPatternPath())

			location := reg.FindStringIndex(path)

			if 0 == len(location) {
				return "", false
			}

			return path[location[1]:], true
		}
	}

	for _, edge := range edges {
		if pathSegment, ok := matcher(edge); ok {
			return edge.node, pathSegment, true
		}
	}

	return nil, "", false
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

func mergeRoutes(routes ...RouteInterface) RouteInterface {

	for i := 1; i <= len(routes); i++ {
		routes[0].AddHandlers(routes[i].GetHandlers())
	}

	return routes[0]
}
