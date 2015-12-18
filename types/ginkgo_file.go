package types

import "fmt"

type GinkgoFile struct {
	// Variables
	BoundVariables []Variable

	// Tree
	containers []*ContainerNode
}

func (g *GinkgoFile) FindNodeById(id string) (Node, error) {
	for _, c := range g.containers {
		if c.Id() == id {
			return c, nil
		}

		if n := c.FindNodeById(id); n != nil {
			return n, nil
		}
	}

	return nil, fmt.Errorf("Node '%s' was not found!", id)
}

func (g *GinkgoFile) BFSIds() []string {
	ids := []string{}

	for _, c := range g.containers {
		ids = append(ids, c.Id())
	}

	for _, c := range g.containers {
		ids = append(ids, c.BFSIds()...)
	}

	return ids
}

func (g *GinkgoFile) AddContainer(c *ContainerNode) {
	g.containers = append(g.containers, c)

	c.SetParent(nil)
}
