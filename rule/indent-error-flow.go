package rule

import (
	"github.com/mgechev/revive/internal/ifelse"
	"github.com/mgechev/revive/lint"
)

// IndentErrorFlowRule lints given else constructs.
type IndentErrorFlowRule struct{}

// Apply applies the rule to given file.
func (e *IndentErrorFlowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	return ifelse.Apply(e, file.AST, ifelse.TargetElse)
}

// Name returns the rule name.
func (*IndentErrorFlowRule) Name() string {
	return "indent-error-flow"
}

func (e *IndentErrorFlowRule) CheckIfElse(chain ifelse.Chain) (failMsg string) {
	if !chain.IfTerminator.DeviatesControlFlow() {
		// this rule only applies if the if-block deviates control flow
		return
	}

	if chain.HasPriorNonReturn {
		// if we de-indent the "else" block then a previous branch
		// might flow into it, affecting program behaviour
		return
	}

	if !chain.IfTerminator.IsReturn() {
		// avoid overlapping with superfluous-else
		return
	}

	return "if block ends with a return statement, so drop this else and outdent its block"
}
