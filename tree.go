package tmux

import (
	"fmt"
	"regexp"
	"strings"
)

// RouteTreeInterface if like you to implement your own tree version, feel free to do it
type RouteTreeInterface interface {
	UseNode(func() *Node)
	Insert(RouteInterface) RouteInterface
	Find(*Node, Method, string) (RouteInterface, map[string]string)
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
	currentNode := t.root
	pathSegments := t.pathSegments(newRoute.GetPattern())
	currentSeg := pathSegments[0]
	currentSegTyp := t.checkNodeType(currentSeg)
	var nextSeg string
	next := true
	if len(pathSegments) > 1 {
		pathSegments = pathSegments[1:]
	} else {
		pathSegments = make([]string, 0, 0)
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
			currentSegTyp = t.checkNodeType(currentSeg)
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
			currentSegTyp := t.checkNodeType(currentSeg)
			n.seg = currentSeg
			currentNode.nodes[currentSegTyp] = append(currentNode.nodes[currentSegTyp], n)
			currentNode = n

			for _, seg := range pathSegments {
				currentSegTyp := t.checkNodeType(seg)
				n := NewNode()
				n.seg = seg
				currentNode.nodes[currentSegTyp] = append(currentNode.nodes[currentSegTyp], n)
				currentNode = n
			}

			currentNode.leaf = newRoute
			return newRoute
		}
	}

	return nil
}

func (t *Tree) checkNodeType(seg string) nodeType {
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

// Find is used to lookup a specific key, returning
// the value and if it was found
func (t *Tree) Find(root *Node, method Method, path string) (RouteInterface, map[string]string) {
	pathSegments := t.pathSegments(path)
	currentSeg := pathSegments[0]
	currentNode := t.root
	copyPathSegments := pathSegments

	if len(pathSegments) > 1 {
		pathSegments = pathSegments[1:]
	} else {
		pathSegments = make([]string, 0, 0)
	}

	for {
		// Node has none sub nodes
		if len(currentNode.nodes[0]) == 0 && len(currentNode.nodes[1]) == 0 && len(currentNode.nodes[2]) == 0 {
			break
		}

	outerLoop:
		for _, typ := range []nodeType{regexNode, staticNode, paramNode} {
			for _, n := range currentNode.nodes[typ] {

				var matched bool
				if regexNode == nodeType(typ) {
					if match, err := regexp.MatchString(n.seg[1:], currentSeg); err == nil && match {
						matched = true
					}
				} else if paramNode == nodeType(typ) && n.seg == ":string" {
					if match, err := regexp.MatchString("([a-zA-Z]{1,})", currentSeg); err == nil && match {
						matched = true
					}
				} else if paramNode == nodeType(typ) && n.seg == ":number" {
					if match, err := regexp.MatchString("([0-9]{1,})", currentSeg); err == nil && match {
						matched = true
					}
				} else if staticNode == nodeType(typ) {
					if n.seg == currentSeg {
						matched = true
					}
				}

				if matched {
					if len(pathSegments) == 0 {
						param := map[string]string{}
						for key, seg := range copyPathSegments {
							param[fmt.Sprintf("seg%d", key)] = seg
						}
						return n.leaf, param
					} else if len(pathSegments) > 1 {
						currentSeg = pathSegments[0]
						pathSegments = pathSegments[1:]
						currentNode = n
					} else if len(pathSegments) == 1 {
						currentSeg = pathSegments[0]
						pathSegments = make([]string, 0, 0)
						currentNode = n
					}
					break outerLoop
				}
			}
		}
	}

	return nil, nil
}

func mergeRoutes(routes ...RouteInterface) RouteInterface {

	for i := 1; i <= len(routes)-1; i++ {
		routes[0].AddHandlers(routes[i].GetHandlers())
	}

	return routes[0]
}
