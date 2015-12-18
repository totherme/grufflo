package types

type Variable interface {
	Name() string
}

type Expr interface {
	Merge(c Expr) (Expr, error)
}

type Node interface {
	Id() string

	// Graph
	Parent() Node
	SetParent(Node)
	IsLeaf() bool
	Children() []Node
	ChildIdx(id string) int
	AddChild(Node) error
	DeleteChild(id string) error
	FindNodeById(id string) Node

	// Children movement
	MoveChildLeft(id string)
	MoveChildRight(id string)
	MoveChildTo(id string, idx int)

	// Debug
	BFSIds() []string
}
