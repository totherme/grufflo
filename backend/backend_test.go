package backend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/totherme/grufflo/backend"
	"github.com/totherme/grufflo/backend/fakes"
	"github.com/totherme/grufflo/types"
)

var _ = Describe("Backend", func() {
	var (
		fakeParser *fakes.FakeParser
		logger     lager.Logger

		bcknd *backend.Backend
	)

	BeforeEach(func() {
		fakeParser = new(fakes.FakeParser)

		ginkgoFile := &types.GinkgoFile{
			BoundVariables: []types.Variable{},
		}

		c0 := &types.ContainerNode{
			Identifier: "/0",
			Subject:    "Lifecycle",
			BoundVariables: []types.Variable{
				fakes.FakeVariable("client"),
				fakes.FakeVariable("container"),
			},
			BeforeEach: fakes.FakeExpr("boo"),
			AfterEach:  fakes.FakeExpr("foo"),
		}
		c0.AddChild(&types.SpecNode{
			Identifier: "/0/0",
		})
		c0.AddChild(&types.SpecNode{
			Identifier: "/0/1",
		})
		ginkgoFile.AddContainer(c0)

		c02 := &types.ContainerNode{
			Identifier: "/0/2",
			Subject:    "StreamIn",
			BoundVariables: []types.Variable{
				fakes.FakeVariable("tarStream"),
			},
			BeforeEach: fakes.FakeExpr("bee"),
			AfterEach:  fakes.FakeExpr("fee"),
		}
		c02.AddChild(&types.SpecNode{
			Identifier: "/0/2/0",
			Subject:    "Streams files as root",
		})
		c02.AddChild(&types.SpecNode{
			Identifier: "/0/2/1",
			Subject:    "Streams files as alice",
			FreeVariables: []types.Variable{
				fakes.FakeVariable("tarStream"),
				fakes.FakeVariable("client"),
				fakes.FakeVariable("container"),
			},
		})
		c0.AddChild(c02)

		c03 := &types.ContainerNode{
			Identifier: "/0/3",
			Subject:    "NetIn",
			BoundVariables: []types.Variable{
				fakes.FakeVariable("ip"),
				fakes.FakeVariable("port"),
			},
			BeforeEach: fakes.FakeExpr("baa"),
			AfterEach:  fakes.FakeExpr("taa"),
		}
		c03.AddChild(&types.SpecNode{
			Identifier: "/0/3/0",
			Subject:    "Listens to a port",
			FreeVariables: []types.Variable{
				fakes.FakeVariable("port"),
				fakes.FakeVariable("client"),
				fakes.FakeVariable("container"),
			},
		})
		c03.AddChild(&types.SpecNode{
			Identifier: "/0/3/1",
			Subject:    "Listens to a train station",
			FreeVariables: []types.Variable{
				fakes.FakeVariable("port"),
				fakes.FakeVariable("ip"),
				fakes.FakeVariable("client"),
				fakes.FakeVariable("container"),
			},
		})
		c0.AddChild(c03)

		fakeParser.ParseReturns(ginkgoFile, nil)

		logger = lagertest.NewTestLogger("test")

		bcknd = &backend.Backend{
			Parser: fakeParser,
			Logger: logger,
		}

		Expect(bcknd.Start()).To(Succeed())
	})

	Describe("MoveOut", func() {
		It("should reorder nodes correctly", func() {
			Expect(bcknd.MoveOut("/0/3/0")).To(Succeed())

			ids := bcknd.GinkgoFile().BFSIds()
			Expect(ids).Should(Equal([]string{
				"/0", "/0/0", "/0/1", "/0/2", "/0/3/0", "/0/3", "/0/2/0", "/0/2/1",
				"/0/3/1",
			}))
		})

		It("should update bound variables of new context", func() {
			Expect(bcknd.MoveOut("/0/3/0")).To(Succeed())

			ginkgoFile := bcknd.GinkgoFile()

			n, err := ginkgoFile.FindNodeById("/0")
			Expect(err).NotTo(HaveOccurred())
			newParent := n.(*types.ContainerNode)

			Expect(newParent.BoundVariables).To(Equal([]types.Variable{
				fakes.FakeVariable("client"),
				fakes.FakeVariable("container"),
				fakes.FakeVariable("ip"),
				fakes.FakeVariable("port"),
			}))
			Expect(newParent.BeforeEach).To(Equal(fakes.FakeExpr("boo\n\nbaa")))
			Expect(newParent.AfterEach).To(Equal(fakes.FakeExpr("foo\n\ntaa")))
		})
	})
})
