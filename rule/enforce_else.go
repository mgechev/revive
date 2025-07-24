package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// EnforceElseRule enforces else branches in if... else if chains
type EnforceElseRule struct{}

// Apply applies the rule to given file.
func (r *EnforceElseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintEnforceElseRule{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*EnforceElseRule) Name() string {
	return "enforce-else"
}

type lintEnforceElseRule struct {
	onFailure func(lint.Failure)
	chain     []ast.Node
}

func (w *lintEnforceElseRule) addBranchToChain(branch ast.Node) {
	if w.chain == nil {
		w.chain = []ast.Node{}
	}

	w.chain = append(w.chain, branch)
}

func (w *lintEnforceElseRule) inIfChain() bool {
	return len(w.chain) > 0
}

func (w *lintEnforceElseRule) resetChain() {
	w.chain = []ast.Node{}
}

func (w *lintEnforceElseRule) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	w.walkBranch(ifStmt.Body)

	switch {
	case ifStmt.Else == nil:
		w.addBranchToChain(ifStmt.Body)
		mustReport := w.inIfChain() && !w.allBranchesEndWithJumpStmt(w.chain)
		if mustReport {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ifStmt,
				Category:   lint.FailureCategoryMaintenance,
				Failure:    `"if ... else if" constructs should end with "else" clauses`,
			})
		}

	default: // there is an else branch
		if chainedIf, ok := ifStmt.Else.(*ast.IfStmt); ok {
			w.addBranchToChain(ifStmt.Body)
			w.Visit(chainedIf)
		} else {
			w.walkBranch(ifStmt.Else)
		}
	}

	w.resetChain()
	return nil // if branches already analyzed
}

func (w *lintEnforceElseRule) allBranchesEndWithJumpStmt(branches []ast.Node) bool {
	for _, branch := range branches {
		block, ok := branch.(*ast.BlockStmt)
		if !ok || len(block.List) == 0 {
			return false
		}

		lastStmt := block.List[len(block.List)-1]

		if _, ok := lastStmt.(*ast.ReturnStmt); ok {
			continue
		}

		if jump, ok := lastStmt.(*ast.BranchStmt); ok && jump.Tok == token.BREAK {
			continue
		}

		return false
	}

	return true
}

func (w *lintEnforceElseRule) walkBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	walker := &lintEnforceElseRule{
		onFailure: w.onFailure,
	}

	ast.Walk(walker, branch)
}
