package backend

import (
	"github.com/pivotal-golang/lager"
	"github.com/totherme/grufflo/types"
)

//go:generate counterfeiter . Parser
type Parser interface {
	Parse(filePath string) (*types.GinkgoFile, error)
}

type Backend struct {
	Parser Parser

	FilePath string

	ginkgoFile *types.GinkgoFile

	Logger lager.Logger
}

func (b *Backend) Start() error {
	var err error

	b.ginkgoFile, err = b.Parser.Parse(b.FilePath)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) MoveOut(id string) error {
	log := b.Logger.Session("move-out", lager.Data{"NodeId": id})

	node, err := b.ginkgoFile.FindNodeById(id)
	if err != nil {
		log.Error("find-node-by-id", err)
		return err
	}

	parent := node.Parent()
	if parent == nil {
		log.Debug("is-top-level-container")
		return nil
	}

	if err := parent.DeleteChild(id); err != nil {
		log.Error("delete-child", err)
		return err
	}

	grandParent := parent.Parent()
	if grandParent == nil && node.IsLeaf() {
		log.Debug("is-top-level-spec")
		return nil
	}

	if grandParent == nil {
		log.Debug("adds-to-ginkgo-file")
		b.ginkgoFile.AddContainer(node.(*types.ContainerNode))
		return nil
	}

	log.Debug("adds-to-grandparent-container")
	grandParent.AddChild(node)

	parentsIdx := grandParent.ChildIdx(parent.Id())
	grandParent.MoveChildTo(node.Id(), parentsIdx)

	return nil
}

func (b *Backend) GinkgoFile() *types.GinkgoFile {
	return b.ginkgoFile
}
