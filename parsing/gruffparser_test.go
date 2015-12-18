package parsing_test

import (
	"github.com/totherme/grufflo/backend"
	"github.com/totherme/grufflo/parsing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gruffparser", func() {
	Describe("the free variables thing", func() {
		var (
			prog  string
			fvSet []string
		)

		JustBeforeEach(func() {
			var err error
			fvSet, err = parsing.GetFVs(prog)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when prog is a simple lambda", func() {
			BeforeEach(func() {
				prog = "func(x,y int) int { return x + y + z }"
			})
			It("should be able to find the free variables of prog", func() {
				Expect(fvSet).To(Equal([]string{"z"}))
			})
		})

		Context("when prog has an assignment to a fresh variable in it", func() {
			BeforeEach(func() {
				prog = `func(x int) int {
  			y := 12
  			return x + y + z
  		}`
			})
			It("should be able to find the free variables of prog", func() {
				Expect(fvSet).To(Equal([]string{"z"}))
			})
		})

		Context("when prog has an assignment to a fresh variable from a free variable in it", func() {
			BeforeEach(func() {
				prog = `func(x int) int {
  			y := a
  			return x + y + z
  		}`
			})
			It("should be able to find the free variables of prog", func() {
				Expect(len(fvSet)).To(Equal(2))
				Expect(contains(fvSet, "a")).To(BeTrue())
				Expect(contains(fvSet, "z")).To(BeTrue())
			})
		})

		Context("when prog has an assignment to a fresh variable from a binop of a free variable and a bound variable in it", func() {
			BeforeEach(func() {
				prog = `func(x int) int {
  			y := a + x
  			return x + y + z
  		}`
			})
			It("should be able to find the free variables of prog", func() {
				Expect(len(fvSet)).To(Equal(2))
				Expect(contains(fvSet, "a")).To(BeTrue())
				Expect(contains(fvSet, "z")).To(BeTrue())
			})
		})

		Context("when prog has an assignment to a fresh variable from a function call of a free variable in it", func() {
			BeforeEach(func() {
				prog = `func(x int) int {
  			y := 12
  			return x + y + f(z)
  		}`
			})
			It("should be able to find the free variables of prog", func() {
				Expect(len(fvSet)).To(Equal(2))
				Expect(contains(fvSet, "f")).To(BeTrue())
				Expect(contains(fvSet, "z")).To(BeTrue())
			})
		})
	})
	Describe("the actual parsing thingie", func() {
		It("should parse a simple testfile", func() {
			var gp backend.Parser
			gp = parsing.GruffParser{}
			ginkFile, err := gp.Parse("testdata/simplefile.go")
			Expect(err).NotTo(HaveOccurred())
			Expect(ginkFile).NotTo(BeNil())
			kids := ginkFile.Children()
			Expect(kids).NotTo(BeNil())
			Expect(len(kids)).To(Equal(1))
		})
	})
})

func contains(lst []string, s string) bool {
	for _, it := range lst {
		if s == it {
			return true
		}
	}
	return false
}
