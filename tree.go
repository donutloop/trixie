package tmux

import (
	"net/http"
	"strings"
)

type RouteTreeInterface interface {
	UseRoute(func() RouteInterface)
	UseNode(func() NodeInterface)
	Insert(Method, string, http.Handler) RouteInterface
	Find(Method, string) RouteInterface
}

// RadixTree implements a radix tree. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and
// ordered iteration,
type RadixTree struct {
	root             NodeInterface
	nodeConstructor  func() NodeInterface
	routeConstructor func() RouteInterface
}

// NewRadixTree returns an empty Tree
func NewRadixTree(nodeConstructor func() NodeInterface, routeConstructor func() RouteInterface) func() RouteTreeInterface {
	return func() RouteTreeInterface {
		tree := &RadixTree{}
		tree.UseNode(nodeConstructor)
		tree.UseRoute(routeConstructor)
		tree.root = tree.nodeConstructor()
		return tree
	}
}

func (t *RadixTree) UseRoute(constructer func() RouteInterface) {
	t.routeConstructor = constructer
}

func (t *RadixTree) UseNode(constructer func() NodeInterface) {
	t.nodeConstructor = constructer
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

			parent.AddEdge(Edge{
				label: search[0],
				node:  newNode,
			})
			return newLeaf
		}

		// Determine longest prefix of the search key on match
		commonPrefix := longestPrefix(search, currentNode.GetPrefixPath())
		if commonPrefix == len(currentNode.GetPrefixPath()) {
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		childNode := t.nodeConstructor().SetPrefixPath(search[:commonPrefix])

		parent.ReplaceEdge(Edge{
			label: search[0],
			node:  childNode,
		})

		// Restore the existing node
		childNode.AddEdge(Edge{
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

		childNode.AddEdge(Edge{
			label: search[0],
			node:  newEdgeNode,
		})

		return newLeaf
	}
}

// Find is used to lookup a specific key, returning
// the value and if it was found
func (t *RadixTree) Find(method Method, path string) RouteInterface {
	var currentNode NodeInterface
	currentNode = t.root
	for {
		// Check for key exhaution
		if len(path) == 0 {
			if currentNode.IsLeaf() && currentNode.GetLeaf().HasHandler(method) {
				return currentNode.GetLeaf()
			}
			break
		}

		// Look for an edge
		currentNode = currentNode.GetEdge(path[0])
		if currentNode == nil {
			break
		}

		// Consume the search prefix
		if !strings.HasPrefix(path, currentNode.GetPrefixPath()) {
			break
		}

		path = path[len(currentNode.GetPrefixPath()):]
	}
	return nil
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
