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
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*IdenticalBranchesRule) Name() string {
	return "identical-branches"
}

// lintIdenticalBranches implements the root or main AST walker of the rule.
// This walker will activate other walkers depending on the satement under analysis:
// - simple if ... else ...
// - if ... else if ... chain
// - switch (not yet implemented)
type lintIdenticalBranches struct {
	file      *lint.File // only necessary to retrieve the line number of branches
	onFailure func(lint.Failure)
}

func (w *lintIdenticalBranches) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt:
		if w.isIfElsIf(n) {
			walker := &lintIfChainIdenticalBranches{
				onFailure:  w.onFailure,
				file:       w.file,
				rootWalker: w,
			}

			ast.Walk(walker, n)
			return nil // the walker already analyzed inner branches
		}

		w.lintSimpleIf(n)
		return w

	case *ast.SwitchStmt:
		// TODO later
		return w
	default:
		return w
	}
}

func (w *lintIdenticalBranches) isIfElsIf(n *ast.IfStmt) bool {
	if n.Else == nil {
		return false
	}

	_, ok := n.Else.(*ast.IfStmt)

	return ok
}

type lintSimpleIfIdenticalBranches struct {
	onFailure func(lint.Failure)
}

func (w *lintIdenticalBranches) lintSimpleIf(n *ast.IfStmt) {
	if n.Else == nil {
		return
	}

	elseBranch, ok := n.Else.(*ast.BlockStmt)
	if !ok { // if-else-if construction (should never be the case but keep the check for safer refactoring)
		return
	}

	if w.identicalBranches(n.Body, elseBranch) {
		w.newFailure(n, "both branches of the if are identical", 1.0)
	}
}

func (w *lintIdenticalBranches) identicalBranches(body *ast.BlockStmt, elseBranch *ast.BlockStmt) bool {
	if len(body.List) != len(elseBranch.List) {
		return false // branches don't have the same number of statements
	}

	bodyStr := astutils.GoFmt(body)
	elseStr := astutils.GoFmt(elseBranch)

	return bodyStr == elseStr
}

type lintIfChainIdenticalBranches struct {
	file                *lint.File // only necessary to retrieve the line number of branches
	onFailure           func(lint.Failure)
	branches            []ast.Stmt             // hold branches to compare
	rootWalker          *lintIdenticalBranches // the walker to use to recursively analize inner branches
	hasComplexCondition bool                   // indicates if one of the if conditions is "complex"
}

// addBranch adds a branch to the list of branches to be compared.
func (w *lintIfChainIdenticalBranches) addBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	if w.branches == nil {
		w.resetBranches()
	}

	w.branches = append(w.branches, branch)
}

// resetBranches resets (clears) the list of branches to compare.
func (w *lintIfChainIdenticalBranches) resetBranches() {
	w.branches = []ast.Stmt{}
	w.hasComplexCondition = false
}

func (w *lintIfChainIdenticalBranches) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// recursively analyze the then-branch
	w.walkBranch(n.Body)

	if n.Init == nil { // only check if without initialization to avoid false positives
		w.addBranch(n.Body)
	}

	if w.isComplexCondition(n.Cond) {
		w.hasComplexCondition = true
	}

	if n.Else != nil {
		if chainedIf, ok := n.Else.(*ast.IfStmt); ok {
			w.Visit(chainedIf)
		} else {
			w.addBranch(n.Else)
			w.walkBranch(n.Else)
		}
	}

	identicalBranches := w.identicalBranches(w.branches)
	for _, branchPair := range identicalBranches {
		branchLines := w.getStmtLines(branchPair)
		msg := fmt.Sprintf("this if...else if chain has identical branches (lines %v)", branchLines)
		confidence := 1.0
		if w.hasComplexCondition {
			confidence = 0.8
		}
		w.rootWalker.newFailure(w.branches[0], msg, confidence)
	}

	w.resetBranches()
	return nil
}

// getStmtLines yields the start line number of the given statements.
func (w *lintIfChainIdenticalBranches) getStmtLines(stmts []ast.Stmt) []int {
	result := []int{}
	for _, stmt := range stmts {
		pos := w.file.ToPosition(stmt.Pos())
		result = append(result, pos.Line)
	}
	return result
}

// walkBranch analyzes the given branch.
func (w *lintIfChainIdenticalBranches) walkBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	walker := &lintIfChainIdenticalBranches{
		onFailure:  w.onFailure,
		file:       w.file,
		rootWalker: w.rootWalker,
	}

	ast.Walk(walker, branch)
}

// isComplexCondition returns true if the given expression is "complex", false otherwise.
// An expression is considered complex if it has a function call.
func (w *lintIfChainIdenticalBranches) isComplexCondition(expr ast.Expr) bool {
	calls := astutils.PickNodes(expr, func(n ast.Node) bool {
		_, ok := n.(*ast.CallExpr)
		return ok
	})

	return len(calls) > 0
}

// identicalBranches yields pairs of identical branches from the given branches.
func (*lintIfChainIdenticalBranches) identicalBranches(branches []ast.Stmt) [][]ast.Stmt {
	result := [][]ast.Stmt{}
	if len(branches) < 2 {
		return result // only one branch to compare thus we return
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
			result = append(result, []ast.Stmt{match, branch})
		}

		hashes[hash] = branch
	}

	return result
}

func (w *lintIdenticalBranches) newFailure(node ast.Node, msg string, confidence float64) {
	w.onFailure(lint.Failure{
		Confidence: confidence,
		Node:       node,
		Category:   lint.FailureCategoryLogic,
		Failure:    msg,
	})
}
