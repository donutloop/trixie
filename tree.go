package tmux

import (
	"net/http"
	"regexp"
	"strings"
)

// RouteTreeInterface if like you to implement your own tree version, feel free to do it
type RouteTreeInterface interface {
	UseRoute(func() RouteInterface)
	UseNode(func() NodeInterface)
	Insert(Method, string, http.Handler) RouteInterface
	Find(NodeInterface, Method, string) RouteInterface
	GetRoot() NodeInterface
}

// RadixTree implements RouteTreeInterface. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and ordered iteration.
// based on go-radix ideas (github.com/armon/go-radix)
type RadixTree struct {
	root             NodeInterface
	nodeConstructor  func() NodeInterface
	routeConstructor func() RouteInterface
}

// NewRadixTree returns an empty Radix Tree
func NewRadixTree(nodeConstructor func() NodeInterface, routeConstructor func() RouteInterface) func() RouteTreeInterface {
	return func() RouteTreeInterface {
		tree := &RadixTree{}
		tree.UseNode(nodeConstructor)
		tree.UseRoute(routeConstructor)
		tree.root = tree.nodeConstructor()
		return tree
	}
}

// UseRoute that you can use diffrent route versions
// See RouteInterface for more details (route.go)
func (t *RadixTree) UseRoute(constructer func() RouteInterface) {
	t.routeConstructor = constructer
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
func (t *RadixTree) Insert(method Method, pattern string, handler http.Handler) RouteInterface {
	var parent NodeInterface
	currentNode := t.root
	search := pattern

	for {
		// Handle key exhaution
		if len(search) == 0 {
			if currentNode.IsLeaf() {
				currentNode.GetLeaf().AddHandler(method, handler)
				currentNode.GetLeaf().SetPattern(pattern)
				return currentNode.GetLeaf()
			}

			currentNode.SetLeaf(t.routeConstructor())
			currentNode.GetLeaf().AddHandler(method, handler)
			currentNode.GetLeaf().SetPattern(pattern)
			return currentNode.GetLeaf()
		}

		// Look for the edge
		parent = currentNode
		currentNode = currentNode.GetEdge(search[0])

		// No edge, create one
		if currentNode == nil {
			newLeaf := t.routeConstructor().AddHandler(method, handler).SetPattern(pattern)
			newNode := t.nodeConstructor().SetPrefixPath(search).SetLeaf(newLeaf)

			parent.AddEdge(&Edge{
				label: search[0],
				node:  newNode,
			})
			return newLeaf
		}

		// Determine longest prefix of the search key on currentNode
		commonPrefix := longestPrefix(search, currentNode.GetPrefixPath())

		// Check if they share the same prefix when yes overwrite current search and continue to next iteration
		if commonPrefix == len(currentNode.GetPrefixPath()) {
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		childNode := t.nodeConstructor().SetPrefixPath(search[:commonPrefix])

		parent.ReplaceEdge(&Edge{
			label: search[0],
			node:  childNode,
		})

		// Restore the existing node
		childNode.AddEdge(&Edge{
			label: currentNode.GetPrefixPath()[commonPrefix],
			node:  currentNode,
		})
		currentNode.SetPrefixPath(currentNode.GetPrefixPath()[commonPrefix:])

		// Create a new leaf node
		newLeaf := t.routeConstructor().SetPattern(pattern).AddHandler(method, handler)

		// If the new key is a subset, add to to this node
		search = search[commonPrefix:]

		if len(search) == 0 {
			childNode.SetLeaf(newLeaf)
		}

		// Create a new edge for the node
		newEdgeNode := t.nodeConstructor().SetPrefixPath(search).SetLeaf(newLeaf)

		childNode.AddEdge(&Edge{
			label: search[0],
			node:  newEdgeNode,
		})

		return newLeaf
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
			// found a node, return it
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
