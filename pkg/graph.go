package flower

// Types with thanks to go/callgraph/callgraph.go
type Edge struct {
	Caller *Node
	Callee *Node
}

type Node struct {
	Flow *Flow
	ID   int
	In   []*Edge
	Out  []*Edge
}

type Graph struct {
	Root  *Node           // the distinguished root node
	Nodes map[*Flow]*Node // all nodes by function
}
