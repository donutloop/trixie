package trixie

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// RouteTreeInterface if like you to implement your own tree version, feel free to do it
type RouteTreeInterface interface {
	UseNode(func() *Node)
	Insert(RouteInterface) RouteInterface
	Find(*Node, string) (RouteInterface, map[string]string, error)
	GetRoot() *Node
}

type Tree struct {
	root            *Node
	nodeConstructor func() *Node
}

// NewTree returns an empty Radix Tree
func NewTree(nodeConstructor func() *Node) func() RouteTreeInterface {
	return func() RouteTreeInterface {
		tree := &Tree{}
		tree.UseNode(nodeConstructor)
		tree.root = tree.nodeConstructor()
		tree.root.root = true
		return tree
	}
}

// UseNode that you can use diffrent node versions
// See NodeInterface for more details (node.go)
func (t *Tree) UseNode(constructer func() *Node) {
	t.nodeConstructor = constructer
}

func (t *Tree) GetRoot() *Node {
	return t.root
}

func (t *Tree) pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

// Insert is used to add a new entry or update
// an existing entry.
func (t *Tree) Insert(newRoute RouteInterface) RouteInterface {

	if newRoute.GetPattern() == "/" {
		if t.root.leaf == nil {
			t.root.leaf = newRoute
		} else {
			t.root.leaf = mergeRoutes(t.root.leaf, newRoute)
		}
		return newRoute
	}

	currentNode := t.root
	pathSegments := t.pathSegments(newRoute.GetPattern())
	currentSeg := pathSegments[0]
	currentSegTyp := NodeOfType(currentSeg)
	var nextSeg string
	next := true
	if len(pathSegments) > 1 {
		pathSegments = pathSegments[1:]
	} else {
		pathSegments = make([]string, 0, 0)
		next = false
	}

	for {

		if !next {
			if currentNode.root {
				n := NewNode()
				n.seg = currentSeg
				n.leaf = newRoute
				currentNode.nodes[currentSegTyp] = append(currentNode.nodes[currentSegTyp], n)
				return newRoute
			} else if currentSeg == currentNode.seg {

				if currentNode.leaf == nil {
					currentNode.leaf = newRoute
				} else {
					currentNode.leaf = mergeRoutes(currentNode.leaf, newRoute)
				}

				return newRoute
			}
		}

		if nextSeg != "" {
			currentSeg = nextSeg
			currentSegTyp = NodeOfType(currentSeg)
		}

	outerLoop:
		for typ, nodes := range currentNode.nodes {

			if nodeType(typ) != currentSegTyp {
				continue
			}

			for _, n := range nodes {
				if n.seg == currentSeg {
					if len(pathSegments) == 1 {
						nextSeg = pathSegments[0]
						currentNode = n
						pathSegments = make([]string, 0, 0)
					} else if len(pathSegments) > 1 {
						nextSeg = pathSegments[0]
						pathSegments = pathSegments[1:]
						currentNode = n
					} else if len(pathSegments) == 0 {
						currentNode = n
						next = false
					}

					break outerLoop
				}
			}

			n := NewNode()
			currentSegTyp := NodeOfType(currentSeg)
			n.seg = currentSeg
			currentNode.nodes[currentSegTyp] = append(currentNode.nodes[currentSegTyp], n)
			currentNode = n

			for _, seg := range pathSegments {
				currentSegTyp := NodeOfType(seg)
				n := NewNode()
				n.seg = seg
				currentNode.nodes[currentSegTyp] = append(currentNode.nodes[currentSegTyp], n)
				currentNode = n
			}

			currentNode.leaf = newRoute
			return newRoute
		}
	}

}

// Find is used to lookup a specific key, returning
// the value and if it was found
func (t *Tree) Find(root *Node, path string) (RouteInterface, map[string]string, error) {

	if path == "" {
		return nil, nil, errors.New("empty path")
	}

	if path == "/" {
		if t.root.leaf == nil {
			return nil, nil, errors.New("root is not a leaf")
		} else {
			return t.root.leaf, nil, nil
		}
	}

	pathSegments := t.pathSegments(path)
	currentSeg := pathSegments[0]
	currentNode := *t.root
	copyPathSegments := pathSegments

	if len(pathSegments) > 1 {
		pathSegments = pathSegments[1:]
	} else {
		pathSegments = make([]string, 0, 0)
	}

	for {
		if !hasSubNodes(&currentNode) {
			break
		}

	outerLoop:
		for _, typ := range []nodeType{regexNode, staticNode, paramNode} {
			for _, n := range currentNode.nodes[typ] {
				if match(typ, currentSeg, n.seg) {
					if len(pathSegments) == 0 {
						param := map[string]string{}
						for key, seg := range copyPathSegments {
							param[fmt.Sprintf("seg%d", key)] = seg
						}
						return n.leaf, param, nil
					} else if len(pathSegments) > 1 {
						currentSeg = pathSegments[0]
						pathSegments = pathSegments[1:]
						currentNode = *n
					} else if len(pathSegments) == 1 {
						currentSeg = pathSegments[0]
						pathSegments = make([]string, 0, 0)
						currentNode = *n
					}
					break outerLoop
				}
			}
			currentNode.nodes[typ] = nil
		}
	}

	return nil, nil, errors.New("path not found")
}

func hasSubNodes(n *Node) bool {
	if len(n.nodes[0]) == 0 && len(n.nodes[1]) == 0 && len(n.nodes[2]) == 0 {
		return false
	}
	return true
}

func NodeOfType(seg string) nodeType {
	var segTyp nodeType
	if seg == ":string" || seg == ":number" {
		segTyp = paramNode
	} else if len(seg) > 0 && string(seg[0]) == "#" {
		segTyp = regexNode
	} else {
		segTyp = staticNode
	}

	return segTyp
}

func match(typ nodeType, currentSeg, seg string) (matched bool) {
	if regexNode == nodeType(typ) {
		if match, err := regexp.MatchString(seg[1:], currentSeg); err == nil && match {
			matched = true
		}
	} else if paramNode == nodeType(typ) && seg == ":string" {
		if match, err := regexp.MatchString("([a-zA-Z]{1,})", currentSeg); err == nil && match {
			matched = true
		}
	} else if paramNode == nodeType(typ) && seg == ":number" {
		if match, err := regexp.MatchString("([0-9]{1,})", currentSeg); err == nil && match {
			matched = true
		}
	} else if staticNode == nodeType(typ) {
		if seg == currentSeg {
			matched = true
		}
	}
	return matched
}

func mergeRoutes(routes ...RouteInterface) RouteInterface {

	for i := 1; i <= len(routes)-1; i++ {
		routes[0].AddHandlers(routes[i].GetHandlers())
	}

	return routes[0]
}
