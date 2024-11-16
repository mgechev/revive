// Package rule implements revive's linting rules.
package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"github.com/mgechev/revive/lint"
)

// Based on https://github.com/fzipp/gocyclo

// CyclomaticRule lints given else constructs.
type CyclomaticRule struct {
	maxComplexity int
	sync.Mutex
}

const defaultMaxCyclomaticComplexity = 10

func (r *CyclomaticRule) configure(arguments lint.Arguments) error {
	r.Lock()
	defer r.Unlock()
	if r.maxComplexity != 0 {
		return nil // already configured
	}

	if len(arguments) < 1 {
		r.maxComplexity = defaultMaxCyclomaticComplexity
		return nil
	}

	complexity, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		return fmt.Errorf("invalid argument for cyclomatic complexity; expected int but got %T", arguments[0])
	}
	r.maxComplexity = int(complexity)
	return nil
}

// Apply applies the rule to given file.
func (r *CyclomaticRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	var failures []lint.Failure
	err := r.configure(arguments)
	if err != nil {
		return failures, err
	}

	fileAst := file.AST

	walker := lintCyclomatic{
		file:       file,
		complexity: r.maxComplexity,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	return failures, nil
}

// Name returns the rule name.
func (*CyclomaticRule) Name() string {
	return "cyclomatic"
}

type lintCyclomatic struct {
	file       *lint.File
	complexity int
	onFailure  func(lint.Failure)
}

func (w lintCyclomatic) Visit(_ ast.Node) ast.Visitor {
	f := w.file
	for _, decl := range f.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		c := complexity(fn)
		if c > w.complexity {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Category:   "maintenance",
				Failure: fmt.Sprintf("function %s has cyclomatic complexity %d (> max enabled %d)",
					funcName(fn), c, w.complexity),
				Node: fn,
			})
		}
	}

	return nil
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	declarationHasReceiver := fn.Recv != nil && fn.Recv.NumFields() > 0
	if declarationHasReceiver {
		typ := fn.Recv.List[0].Type
		return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
	}

	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

// complexity calculates the cyclomatic complexity of a function.
func complexity(fn *ast.FuncDecl) int {
	v := complexityVisitor{}
	ast.Walk(&v, fn)
	return v.Complexity
}

type complexityVisitor struct {
	// Complexity is the cyclomatic complexity
	Complexity int
}

// Visit implements the ast.Visitor interface.
func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
		v.Complexity++
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.Complexity++
		}
	}
	return v
}
