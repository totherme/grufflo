package types

import "errors"

type SpecNode struct {
	Identifier string

	// Ginkgo
	Subject string
	Body    Expr

	// Variables
	FreeVariables []Variable

	// Tree
	parent Node
}

func (s *SpecNode) Id() string {
	return s.Identifier
}

func (s *SpecNode) Parent() Node {
	return s.parent
}

func (s *SpecNode) SetParent(n Node) {
	s.parent = n
}

func (s *SpecNode) IsLeaf() bool {
	return true
}

func (s *SpecNode) Children() []Node {
	return []Node{}
}

func (s *SpecNode) AddChild(n Node) error {
	return errors.New("Spec is a leaf node!")
}

func (s *SpecNode) DeleteChild(_ string) error {
	return errors.New("Spec is a leaf node!")
}

func (s *SpecNode) FindNodeById(_ string) Node {
	// it is always leaf
	return nil
}

func (s *SpecNode) BFSIds() []string {
	// it is always leaf
	return []string{}
}
