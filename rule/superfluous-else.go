package rule

import (
	"fmt"
	"github.com/mgechev/revive/internal/ifelse"
	"github.com/mgechev/revive/lint"
)

// SuperfluousElseRule lints given else constructs.
type SuperfluousElseRule struct{}

// Apply applies the rule to given file.
func (e *SuperfluousElseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	return ifelse.Apply(e, file.AST, ifelse.TargetElse)
}

// Name returns the rule name.
func (*SuperfluousElseRule) Name() string {
	return "superfluous-else"
}

func (e *SuperfluousElseRule) CheckIfElse(chain ifelse.Chain) (failMsg string) {
	if !chain.IfTerminator.DeviatesControlFlow() {
		// this rule only applies if the if-block deviates control flow
		return
	}

	if chain.HasPriorNonReturn {
		// if we de-indent the "else" block then a previous branch
		// might flow into it, affecting program behaviour
		return
	}

	if chain.IfTerminator.IsReturn() {
		// avoid overlapping with indent-error-flow
		return
	}

	return fmt.Sprintf("if block ends with %v, so drop this else and outdent its block", chain.IfTerminator.LongString())
}
