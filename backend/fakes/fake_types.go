package fakes

import (
	"errors"
	"fmt"

	"github.com/totherme/grufflo/types"
)

type FakeExpr string

func (c FakeExpr) Merge(otherExpr types.Expr) (types.Expr, error) {
	otherFakeExpr, ok := otherExpr.(FakeExpr)
	if !ok {
		return nil, errors.New("everybody is fake!")
	}

	var ret FakeExpr = FakeExpr(fmt.Sprintf("%s\n\n%s", c, otherFakeExpr))
	return ret, nil
}

type FakeVariable string

func (v FakeVariable) Name() string {
	return string(v)
}
