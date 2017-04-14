package trixie

type nodeType int

const (
	staticNode nodeType = iota
	paramNode
	regexNode
	nodeTypes
)

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
	root bool

	// leaf is used to store possible leaf
	leaf RouteInterface

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	nodes [nodeTypes][]*Node

	seg string

	param map[string]string
}
