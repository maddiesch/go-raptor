package conditional

import (
	"github.com/maddiesch/go-raptor/statement/generator"
)

type Conditional interface {
	Generate(generator.ArgumentNameProvider) (string, []any)
}
