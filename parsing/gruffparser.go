package parsing

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

	"github.com/totherme/grufflo/types"
)

var idgen int

func GetFVs(code string) ([]string, error) {
	codeAst, err := parser.ParseExpr(code)
	if err != nil {
		return nil, err
	}
	gatherer := newFvGatherer()
	ast.Walk(gatherer, codeAst)
	return gatherer.FVs(), nil
}

func newFvGatherer() *fvGatherer {
	return &fvGatherer{
		freevars:  make(map[string]struct{}),
		boundvars: make(map[string]struct{}),
	}
}

type fvGatherer struct {
	freevars  map[string]struct{}
	boundvars map[string]struct{}
}

func (g *fvGatherer) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	default:
		typ := reflect.TypeOf(n)
		if typ != nil {
			fmt.Printf("node type: %s\n", typ.String())
		}
	case *ast.Ident:
		_, inFree := g.freevars[n.Name]
		_, inBound := g.boundvars[n.Name]
		if !inFree && !inBound {
			g.freevars[n.Name] = struct{}{}
		}
	case *ast.AssignStmt:
		rhs := n.Rhs
		return newBvGatherer(g, rhs)
	case *ast.FuncType:
		return newBvGatherer(g, []ast.Node{})
	}

	return g
}

func (g *fvGatherer) FVs() []string {
	keys := make([]string, len(g.freevars))

	i := 0
	for k := range g.freevars {
		keys[i] = k
		i++
	}

	return keys
}

// The argument breakouts should be of type []ast.Node. I don't understand go
// well enough to make the type checker check for this, without losing the
// ability to pass something of type []ast.Expr, even though ast.Expr _is_ an
// ast.Node
func newBvGatherer(fvg *fvGatherer, breakouts interface{}) *bvGatherer {
	g := bvGatherer{
		fvg:          fvg,
		breakOutList: make(map[ast.Node]struct{}),
	}
	switch breakouts := breakouts.(type) {
	default:
		typ := reflect.TypeOf(breakouts)
		if typ != nil {
			panic(fmt.Sprintf("I just can't coppe. Type of breakouts is: %s\n", typ.String()))
		}
		panic("wtf")
	case []ast.Node:
		for _, b := range breakouts {
			g.breakOutList[b] = struct{}{}
		}
	case []ast.Expr:
		for _, b := range breakouts {
			var bNode ast.Node
			bNode = b
			g.breakOutList[bNode] = struct{}{}
		}
	}
	return &g
}

type bvGatherer struct {
	fvg          *fvGatherer
	breakOutList map[ast.Node]struct{}
}

func (g *bvGatherer) Visit(n ast.Node) ast.Visitor {
	if _, ok := g.breakOutList[n]; ok {
		return g.fvg.Visit(n)
	}

	switch n := n.(type) {
	default:
	case *ast.Ident:
		if _, ok := g.fvg.boundvars[n.Name]; !ok {
			g.fvg.boundvars[n.Name] = struct{}{}
		}
	}

	return g
}

type GruffParser struct{}

func (gp GruffParser) Parse(fp string) (*types.GinkgoFile, error) {
	fset := token.NewFileSet()

	codeAst, err := parser.ParseFile(fset, fp, nil, 0)
	if err != nil {
		return nil, err
	}

	gfp := &ginkFileParser{
		ginkFile: &types.GinkgoFile{},
	}
	ast.Walk(gfp, codeAst)

	return gfp.ginkFile, nil
}

type ginkFileParser struct {
	ginkFile *types.GinkgoFile
}

func (p *ginkFileParser) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	default:
		typ := reflect.TypeOf(n)
		if typ != nil {
			fmt.Printf("GFP: node type: %s\n", typ.String())
		}
	case *ast.CallExpr:
		// This could be a describe block!
		expr := n.Fun
		switch expr := expr.(type) {
		default:
			break
		case *ast.Ident:
			if expr.Name == "Describe" {
				// This IS a describe block!
				fmt.Println("This IS a describe block!")

				describe := &types.ContainerNode{
					Identifier: newID(),
				}
				p.ginkFile.AddContainer(describe)
				return &ginkContainerParser{ginkContainer: describe}
			}
		}
		return p
	}
	return p
}

type ginkContainerParser struct {
	ginkContainer *types.ContainerNode
}

func (p *ginkContainerParser) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	default:
		typ := reflect.TypeOf(n)
		if typ != nil {
			fmt.Printf("GCP: node type: %s\n", typ.String())
		}
	case *ast.BasicLit:
		fmt.Printf("GODP: basic lit with value: %s\n", n.Value)
		p.ginkContainer.Subject = n.Value
	}
	return p
}

func newID() string {
	return fmt.Sprintf("%d", idgen)
}
