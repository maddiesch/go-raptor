package conditional

import (
	"fmt"

	"github.com/maddiesch/go-raptor/statement/generator"
)

func And(left, right Conditional) Conditional {
	return &logicalInfixConditional{left, right, "AND"}
}

func Or(left, right Conditional) Conditional {
	return &logicalInfixConditional{left, right, "OR"}
}

type logicalInfixConditional struct {
	left     Conditional
	right    Conditional
	operator string
}

func (c *logicalInfixConditional) Generate(provider generator.ArgumentNameProvider) (string, []any) {
	if c.left == nil && c.right == nil {
		panic("must provide at least one non-nil conditional")
	}
	if c.left == nil {
		return c.right.Generate(provider)
	}
	if c.right == nil {
		return c.left.Generate(provider)
	}
	var args []any

	left, lArgs := c.left.Generate(provider)
	args = append(args, lArgs...)

	right, rArgs := c.right.Generate(provider)
	args = append(args, rArgs...)

	return fmt.Sprintf("(%s %s %s)", left, c.operator, right), args
}
