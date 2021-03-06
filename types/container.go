package types

import "fmt"

type ContainerNode struct {
	Identifier string

	// Ginkgo
	Subject    string
	BeforeEach Expr
	AfterEach  Expr

	// Variables
	BoundVariables []Variable

	// Tree
	nodes  []Node
	parent Node
}

func (c *ContainerNode) Id() string {
	return c.Identifier
}

func (c *ContainerNode) Parent() Node {
	return c.parent
}

func (c *ContainerNode) SetParent(n Node) {
	c.parent = n
}

func (c *ContainerNode) IsLeaf() bool {
	return false
}

func (c *ContainerNode) Children() []Node {
	return c.nodes
}

func (c *ContainerNode) ChildIdx(id string) int {
	for i := 0; i < len(c.nodes); i++ {
		if c.nodes[i].Id() == id {
			return i
		}
	}

	return -1
}

func (c *ContainerNode) AddChild(n Node) error {
	n.SetParent(c)

	c.nodes = append(c.nodes, n)

	return nil
}

func (c *ContainerNode) DeleteChild(id string) error {
	i := c.ChildIdx(id)

	if i == len(c.nodes) {
		return fmt.Errorf("Node '%s' was not found!", id)
	}

	head := []Node{}
	if i > 0 {
		head = c.nodes[:i]
	}

	tail := []Node{}
	if i < len(c.nodes)-1 {
		tail = c.nodes[i+1:]
	}

	c.nodes = append(head, tail...)

	return nil
}

func (c *ContainerNode) MoveChildLeft(id string) {
	idx := c.ChildIdx(id)
	if idx == -1 || idx == 0 {
		return
	}

	c.moveChildFromTo(idx, idx-1)
}

func (c *ContainerNode) MoveChildRight(id string) {
	idx := c.ChildIdx(id)
	if idx == -1 || idx == len(c.nodes)-1 {
		return
	}

	c.moveChildFromTo(idx, idx+1)
}

func (c *ContainerNode) MoveChildTo(id string, idx int) {
	if idx >= len(c.nodes) {
		return
	}

	fromIdx := c.ChildIdx(id)
	if fromIdx == -1 || fromIdx == idx {
		return
	}

	c.moveChildFromTo(fromIdx, idx)
}

func (c *ContainerNode) FindNodeById(id string) Node {
	containers := []*ContainerNode{}

	for _, n := range c.nodes {
		if id == n.Id() {
			return n
		}

		if c, ok := n.(*ContainerNode); ok {
			containers = append(containers, c)
		}
	}

	for _, c := range containers {
		if n := c.FindNodeById(id); n != nil {
			return n
		}
	}

	return nil
}

func (c *ContainerNode) BFSIds() []string {
	ids := []string{}
	containers := []*ContainerNode{}

	for _, n := range c.nodes {
		ids = append(ids, n.Id())

		if c, ok := n.(*ContainerNode); ok {
			containers = append(containers, c)
		}
	}

	for _, c := range containers {
		cIds := c.BFSIds()
		ids = append(ids, cIds...)
	}

	return ids
}

func (c *ContainerNode) moveChildFromTo(fromIdx, toIdx int) {
	newNodes := make([]Node, len(c.nodes))

	newNodes[toIdx] = c.nodes[fromIdx]

	var a, b int
	if fromIdx < toIdx {
		a = fromIdx
		b = toIdx

		copy(newNodes[fromIdx:toIdx], c.nodes[fromIdx+1:toIdx+1])
	} else {
		a = toIdx
		b = fromIdx

		copy(newNodes[toIdx+1:fromIdx+1], c.nodes[toIdx:fromIdx])
	}

	if a > 0 {
		copy(newNodes[:a], c.nodes[:a])
	}

	if b < len(c.nodes)-1 {
		copy(newNodes[b+1:], c.nodes[b+1:])
	}

	c.nodes = newNodes
}
