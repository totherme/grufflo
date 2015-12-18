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
			BeforeEach: fakes.FakeExpr("I am setting up client and container"),
			AfterEach:  fakes.FakeExpr("I am tearing down client and container"),
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
			BeforeEach: fakes.FakeExpr("I am setting up tarStream"),
			AfterEach:  fakes.FakeExpr("I am tearing down tarStream"),
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
			BeforeEach: fakes.FakeExpr("I am setting up ip and port"),
			AfterEach:  fakes.FakeExpr("I am tearing down ip and port"),
		}
		c03.AddChild(&types.SpecNode{
			Identifier: "/0/3/0",
			Subject:    "Listens to a port",
		})
		c03.AddChild(&types.SpecNode{
			Identifier: "/0/3/1",
			Subject:    "Listens to a train station",
		})
		c0.AddChild(c03)

		fakeParser.ParseReturns(ginkgoFile, nil)

		logger = lagertest.NewTestLogger("test")

		bcknd = &backend.Backend{
			Parser: fakeParser,
			Logger: logger,
		}
	})

	// Notes: we could run a tree-search algorithm and assert on the expected
	//	order of the ids.
	It("should work", func() {
		Expect(bcknd.Start()).To(Succeed())

		Expect(bcknd.MoveOut("/0/3/0")).To(Succeed())
		ids := bcknd.GinkgoFile().BFSIds()
		Expect(ids).Should(Equal([]string{
			"/0", "/0/0", "/0/1", "/0/2", "/0/3/0", "/0/3", "/0/2/0", "/0/2/1",
			"/0/3/1",
		}))
	})
})
