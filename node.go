package trixie

type nodeType int

// All kind of Nodes
const (
	staticNode nodeType = iota
	paramNode
	regexNode
	nodeTypes
)

// NewNode creates a Node instance and setup place for sub nodes
func NewNode() *Node {
	nodes := [nodeTypes][]*Node{}
	nodes[staticNode] = make([]*Node, 0, 0)
	nodes[paramNode] = make([]*Node, 0, 0)
	nodes[regexNode] = make([]*Node, 0, 0)
	return &Node{
		nodes: nodes,
	}
}

type Node struct {
	// is current node the root node of a tree
	root bool

	// leaf (Route instance) is used to store possible leaf
	leaf RouteInterface

	// Node should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	nodes [nodeTypes][]*Node

	// Segment of an path
	seg string
}
