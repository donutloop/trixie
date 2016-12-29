package tmux

import (
	"errors"
	"sort"
	"strings"
)

type NodeInterface interface {
	SetPrefixPath(string) NodeInterface
	GetPrefixPath() string
	SetPrefixPatternPath(string) NodeInterface
	GetPrefixPatternPath() string
	GetEdges() [edgeTypes]Edges
	AddEdge(*Edge)
	GetEdge(label byte) NodeInterface
	ReplaceEdge(e *Edge) error
	IsLeaf() bool
	GetLeaf() RouteInterface
	SetLeaf(RouteInterface) NodeInterface
}

type egdeType uint8

const (
	staticNode egdeType = iota
	paramNode
	regexNode
	edgeTypes
)

func NewNode() NodeInterface {
	return &Node{
		edges: [edgeTypes]Edges{},
	}
}

type Node struct {
	// leaf is used to store possible leaf
	leaf RouteInterface

	// prefix is the common prefix we ignore
	prefix string

	//prefixPattern is the common prefix in regex format
	prefixPattern string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges [edgeTypes]Edges
}

func (n *Node) IsLeaf() bool {
	return n.leaf != nil
}

func (n *Node) SetPrefixPath(prefix string) NodeInterface {
	n.prefix = prefix
	return n
}

func (n *Node) GetPrefixPath() string {
	return n.prefix
}

func (n *Node) SetPrefixPatternPath(prefix string) NodeInterface {
	n.prefixPattern = prefix
	return n
}

func (n *Node) GetPrefixPatternPath() string {
	return n.prefixPattern
}

func (n *Node) SetLeaf(leaf RouteInterface) NodeInterface {
	n.leaf = leaf
	return n
}

func (n *Node) GetLeaf() RouteInterface {
	return n.leaf
}

func (n *Node) AddEdge(e *Edge) {

	n.AddType(e)
	n.PopulatePattern(e)
	n.edges[e.typ] = append(n.edges[e.typ], e)
	n.edges[e.typ].Sort()
}

func (n *Node) AddType(e *Edge) {
	prefixPath := e.node.GetPrefixPath()

	if strings.Contains(prefixPath, ":") {
		e.typ = paramNode
	} else {
		e.typ = staticNode
	}
}

func (n *Node) PopulatePattern(e *Edge) {
	if e.typ == paramNode {
		e.node.SetPrefixPatternPath("^" + strings.Replace(e.node.GetPrefixPath(), ":number", "([0-9]{1,})", -1))
	}
}

func (n *Node) GetEdges() [edgeTypes]Edges {
	return n.edges
}

func (n *Node) GetEdge(label byte) NodeInterface {

	for _, edges := range n.edges {
		for _, edge := range edges {
			if edge.label == label {
				return edge.node
			}
		}
	}

	return nil
}

func (n *Node) ReplaceEdge(e *Edge) error {

	n.AddType(e)

	for i := 0; i < len(n.edges); i++ {
		for j := 0; j < len(n.edges[i]); j++ {
			if n.edges[i][j].label == e.label {
				n.PopulatePattern(e)
				n.edges[i] = append(n.edges[i][:j], n.edges[i][j+1:]...)
				n.edges[e.typ] = append(n.edges[e.typ], e)
				return nil
			}
		}
	}

	return errors.New("replacing missing edge")
}

// Edge is used to represent an edge node
type Edge struct {
	typ   egdeType
	label byte
	node  NodeInterface
}

type Edges []*Edge

func (e Edges) Len() int {
	return len(e)
}

func (e Edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e Edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Edges) Sort() {
	sort.Sort(e)
}
