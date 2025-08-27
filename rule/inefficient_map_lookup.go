package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// InefficientMapLookupRule spots potential inefficient map lookups.
type InefficientMapLookupRule struct{}

// Apply applies the rule to given file.
func (*InefficientMapLookupRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInefficientMapLookup{
		file:      file,
		onFailure: onFailure,
	}

	if err := file.Pkg.TypeCheck(); err != nil {
		return []lint.Failure{
			lint.NewInternalFailure(fmt.Sprintf("Unable to type check file %q: %v", file.Name, err)),
		}
	}

	// Iterate over declarations looking for function declarations
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue // not a function
		}

		if fn.Body == nil {
			continue // external (no-Go) function
		}

		// Analyze the function body
		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*InefficientMapLookupRule) Name() string {
	return "inefficient-map-lookup"
}

type lintInefficientMapLookup struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintInefficientMapLookup) Visit(node ast.Node) ast.Visitor {
	// Only interested in blocks of statements
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w // not a block of statements
	}

	w.analyzeBlock(block)

	return w
}

// analyzeBlock searches AST subtrees with the following form
//
//	for <key> := range <map> {
//		if <key> == <something> {
//		   ...
//		}
func (w *lintInefficientMapLookup) analyzeBlock(b *ast.BlockStmt) {
	for _, stmt := range b.List {
		if !w.isRangeOverMapKey(stmt) {
			continue
		}

		rangeOverMap := stmt.(*ast.RangeStmt)
		key := rangeOverMap.Key.(*ast.Ident)

		// Here we have identified a range over the keys of a map
		// Let's check if the range body is just an if of the form
		// if <key> == <something>
		if len(rangeOverMap.Body.List) != 1 {
			continue // the range body has more than one statement
		}

		bodyStmt := rangeOverMap.Body.List[0]
		ifStmt, ok := bodyStmt.(*ast.IfStmt)
		if !ok {
			continue // the single statement of the body is not an if
		}

		binExp, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok || binExp.Op != token.EQL {
			continue // the if condition is not ==
		}

		if !astutils.IsIdent(binExp.X, key.Name) {
			continue // the if condition is not <key> ==
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       rangeOverMap,
			Category:   lint.FailureCategoryCodeStyle,
			Failure:    "inefficient lookup of map key",
		})
	}
}

func (w *lintInefficientMapLookup) isRangeOverMapKey(stmt ast.Stmt) bool {
	rangeStmt, ok := stmt.(*ast.RangeStmt)
	if !ok {
		return false // not a range
	}

	// Check if we range only on key
	// for key := range ...
	// for key, _ := range ...
	hasValueVariable := rangeStmt.Value != nil && !astutils.IsIdent(rangeStmt.Value, "_")
	if hasValueVariable {
		return false // range over both key and value
	}

	// Check if we range over a map
	t := w.file.Pkg.TypeOf(rangeStmt.X)
	return t != nil && strings.HasPrefix(t.String(), "map[")
}
