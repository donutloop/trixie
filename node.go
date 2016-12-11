package tmux

import "sort"

type NodeInterface interface {
	IsLeaf() bool
	AddEdge(Edge)
	GetEdge(label byte) NodeInterface
	ReplaceEdge(e Edge)
	SetLeaf(RouteInterface) NodeInterface
	SetPrefixPath(string) NodeInterface
	GetLeaf() RouteInterface
	GetPrefixPath() string
}

func NewNode() NodeInterface {
	return &Node{}
}

type Node struct {
	// leaf is used to store possible leaf
	leaf RouteInterface

	// prefix is the common prefix we ignore
	prefix string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges Edges
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

func (n *Node) SetLeaf(leaf RouteInterface) NodeInterface {
	n.leaf = leaf
	return n
}

func (n *Node) GetLeaf() RouteInterface {
	return n.leaf
}

func (n *Node) AddEdge(e Edge) {
	n.edges = append(n.edges, e)
	n.edges.Sort()
}

func (n *Node) GetEdge(label byte) NodeInterface {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		return n.edges[idx].node
	}
	return nil
}

func (n *Node) ReplaceEdge(e Edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})
	if idx < num && n.edges[idx].label == e.label {
		n.edges[idx].node = e.node
		return
	}
	panic("replacing missing edge")
}

// edge is used to represent an edge node
type Edge struct {
	label byte
	node  NodeInterface
}

type Edges []Edge

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
