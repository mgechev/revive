package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// IdenticalIfElseIfBranchesRule warns if...else if chains with identical branches.
type IdenticalIfElseIfBranchesRule struct{}

// Apply applies the rule to given file.
func (*IdenticalIfElseIfBranchesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &rootWalkerIfElseIfIdenticalBranches{file: file, onFailure: onFailure}
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
func (*IdenticalIfElseIfBranchesRule) Name() string {
	return "identical-ifelseif-branches"
}

type rootWalkerIfElseIfIdenticalBranches struct {
	file      *lint.File // only necessary to retrieve the line number of branches
	onFailure func(lint.Failure)
}

func (w *rootWalkerIfElseIfIdenticalBranches) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt:
		if w.isIfElseIf(n) {
			walker := &lintIfChainIdenticalBranches{
				onFailure:  w.onFailure,
				file:       w.file,
				rootWalker: w,
			}

			ast.Walk(walker, n)
			return nil // the walker already analyzed inner branches
		}

		return w
	default:
		return w
	}
}

func (w *rootWalkerIfElseIfIdenticalBranches) lintSwitch(switchStmt *ast.SwitchStmt) {
	doesFallthrough := func(stmts []ast.Stmt) bool {
		if len(stmts) == 0 {
			return false
		}

		ft, ok := stmts[len(stmts)-1].(*ast.BranchStmt)
		return ok && ft.Tok == token.FALLTHROUGH
	}

	hashes := map[string]int{} // map hash(branch code) -> branch line
	for _, cc := range switchStmt.Body.List {
		caseClause := cc.(*ast.CaseClause)
		if doesFallthrough(caseClause.Body) {
			continue // skip fallthrough branches
		}
		branch := &ast.BlockStmt{
			List: caseClause.Body,
		}
		hash := astutils.NodeHash(branch)
		branchLine := w.file.ToPosition(caseClause.Pos()).Line
		if matchLine, ok := hashes[hash]; ok {
			w.newFailure(
				switchStmt,
				fmt.Sprintf(`"switch" with identical branches (lines %d and %d)`, matchLine, branchLine),
				1.0,
			)
		}

		hashes[hash] = branchLine
		w.walkBranch(branch)
	}
}

// walkBranch analyzes the given branch.
func (w *rootWalkerIfElseIfIdenticalBranches) walkBranch(branch ast.Stmt) {
	if branch == nil {
		return
	}

	walker := &rootWalkerIfElseIfIdenticalBranches{
		onFailure: w.onFailure,
		file:      w.file,
	}

	ast.Walk(walker, branch)
}

func (*rootWalkerIfElseIfIdenticalBranches) isIfElseIf(node *ast.IfStmt) bool {
	_, ok := node.Else.(*ast.IfStmt)
	return ok
}

func (w *rootWalkerIfElseIfIdenticalBranches) lintSimpleIf(n *ast.IfStmt) {
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

func (*rootWalkerIfElseIfIdenticalBranches) identicalBranches(body, elseBranch *ast.BlockStmt) bool {
	if len(body.List) != len(elseBranch.List) {
		return false // branches don't have the same number of statements
	}

	bodyStr := astutils.GoFmt(body)
	elseStr := astutils.GoFmt(elseBranch)

	return bodyStr == elseStr
}

func (w *rootWalkerIfElseIfIdenticalBranches) newFailure(node ast.Node, msg string, confidence float64) {
	w.onFailure(lint.Failure{
		Confidence: confidence,
		Node:       node,
		Category:   lint.FailureCategoryLogic,
		Failure:    msg,
	})
}

type lintIfChainIdenticalBranches struct {
	file                *lint.File // only necessary to retrieve the line number of branches
	onFailure           func(lint.Failure)
	branches            []ast.Stmt                           // hold branches to compare
	rootWalker          *rootWalkerIfElseIfIdenticalBranches // the walker to use to recursively analize inner branches
	hasComplexCondition bool                                 // indicates if one of the if conditions is "complex"
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
	w.rootWalker.walkBranch(n.Body)

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
			w.rootWalker.walkBranch(n.Else)
		}
	}

	identicalBranches := w.identicalBranches(w.branches)
	for _, branchPair := range identicalBranches {
		msg := fmt.Sprintf(`"if...else if" chain with identical branches (lines %d and %d)`, branchPair[0], branchPair[1])
		confidence := 1.0
		if w.hasComplexCondition {
			confidence = 0.8
		}
		w.rootWalker.newFailure(w.branches[0], msg, confidence)
	}

	w.resetBranches()
	return nil
}

// isComplexCondition returns true if the given expression is "complex", false otherwise.
// An expression is considered complex if it has a function call.
func (*lintIfChainIdenticalBranches) isComplexCondition(expr ast.Expr) bool {
	calls := astutils.PickNodes(expr, func(n ast.Node) bool {
		_, ok := n.(*ast.CallExpr)
		return ok
	})

	return len(calls) > 0
}

// identicalBranches yields pairs of (line numbers) of identical branches from the given branches.
func (w *lintIfChainIdenticalBranches) identicalBranches(branches []ast.Stmt) [][]int {
	result := [][]int{}
	if len(branches) < 2 {
		return result // only one branch to compare thus we return
	}

	hashes := map[string]int{} // branch code hash -> branch line
	for _, branch := range branches {
		hash := astutils.NodeHash(branch)
		branchLine := w.file.ToPosition(branch.Pos()).Line
		if match, ok := hashes[hash]; ok {
			result = append(result, []int{match, branchLine})
		}

		hashes[hash] = branchLine
	}

	return result
}
