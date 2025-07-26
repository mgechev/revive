package rule

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// IdenticalBranchesRule warns on constant logical expressions.
type IdenticalBranchesRule struct{}

// Apply applies the rule to given file.
func (*IdenticalBranchesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintIdenticalBranches{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*IdenticalBranchesRule) Name() string {
	return "identical-branches"
}

type lintIdenticalBranches struct {
	file      *lint.File
	onFailure func(lint.Failure)
	branches  []ast.Stmt
}

func (w *lintIdenticalBranches) addBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	if w.branches == nil {
		w.resetBranches()
	}

	w.branches = append(w.branches, branch)
}

func (w *lintIdenticalBranches) resetBranches() {
	w.branches = []ast.Stmt{}
}

func (w *lintIdenticalBranches) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// recursevely analyze the then-branch
	w.walkBranch(n.Body)
	w.addBranch(n.Body)

	if n.Else != nil {
		if chainedIf, ok := n.Else.(*ast.IfStmt); ok {
			w.Visit(chainedIf)
		} else {
			w.addBranch(n.Else)
			w.walkBranch(n.Else)
		}
	}

	if matching := w.identicalBranches(w.branches); matching != nil {
		msg := "both branches of the if are identical"
		if len(w.branches) > 2 {
			branchLines := w.getStmtLines(matching)
			msg = fmt.Sprintf("this if...else if chain has identical branches (lines %v)", branchLines)
		}

		w.newFailure(w.branches[0], msg)
	}

	w.resetBranches()
	return nil
}

func (w *lintIdenticalBranches) getStmtLines(stmts []ast.Stmt) []int {
	result := []int{}
	for _, stmt := range stmts {
		pos := w.file.ToPosition(stmt.Pos())
		result = append(result, pos.Line)
	}
	return result
}

func (w *lintIdenticalBranches) walkBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	walker := &lintIdenticalBranches{
		onFailure: w.onFailure,
	}

	ast.Walk(walker, branch)
}

func (*lintIdenticalBranches) identicalBranches(branches []ast.Stmt) []ast.Stmt {
	if len(branches) < 2 {
		return nil // only one branch to compare thus we return
	}

	hasher := func(in string) string {
		binHash := md5.Sum([]byte(in))
		return hex.EncodeToString(binHash[:])
	}

	hashes := map[string]ast.Stmt{}
	for _, branch := range branches {
		str := astutils.GoFmt(branch)
		hash := hasher(str)

		if match, ok := hashes[hash]; ok {
			return []ast.Stmt{match, branch}
		}

		hashes[hash] = branch
	}

	return nil
}

func (w *lintIdenticalBranches) newFailure(node ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryLogic,
		Failure:    msg,
	})
}
