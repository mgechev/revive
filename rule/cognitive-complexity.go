package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// CognitiveComplexityRule lints given else constructs.
type CognitiveComplexityRule struct{}

// Apply applies the rule to given file.
func (r *CognitiveComplexityRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	const expectedArgumentsCount = 1
	if len(arguments) < expectedArgumentsCount {
		panic(fmt.Sprintf("not enough arguments for cognitive-complexity, expected %d, got %d", expectedArgumentsCount, len(arguments)))
	}
	complexity, ok := arguments[0].(int64)
	if !ok {
		panic(fmt.Sprintf("invalid argument type for cognitive-complexity, expected int64, got %T", arguments[0]))
	}

	linter := cognitiveComplexityLinter{
		file:          file,
		maxComplexity: int(complexity),
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	linter.lint()

	return failures
}

// Name returns the rule name.
func (r *CognitiveComplexityRule) Name() string {
	return "cognitive-complexity"
}

type cognitiveComplexityLinter struct {
	file          *lint.File
	maxComplexity int
	onFailure     func(lint.Failure)
}

func (w cognitiveComplexityLinter) lint() ast.Visitor {
	f := w.file
	for _, decl := range f.AST.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			v := cognitiveComplexityVisitor{}
			c := v.subTreeComplexity(fn.Body)
			if c > w.maxComplexity {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Category:   "maintenance",
					Failure:    fmt.Sprintf("function %s has cognitive complexity %d (> max enabled %d)", funcName(fn), c, w.maxComplexity),
					Node:       fn,
				})
			}
		}
	}

	return nil
}

type cognitiveComplexityVisitor struct {
	complexity   int
	nestingLevel int
}

// subTreeComplexity calculates the cognitive complexity of an AST-subtree.
func (v cognitiveComplexityVisitor) subTreeComplexity(n ast.Node) int {
	ast.Walk(&v, n)
	return v.complexity
}

// Visit implements the ast.Visitor interface.
func (v *cognitiveComplexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.IfStmt:
		targets := []ast.Node{n.Cond, n.Body, n.Else}
		v.walk(1, targets...)
		return nil
	case *ast.ForStmt:
		targets := []ast.Node{n.Cond, n.Body}
		v.walk(1, targets...)
		return nil
	case *ast.RangeStmt:
		v.walk(1, n.Body)
		return nil
	case *ast.SelectStmt:
		v.walk(1, n.Body)
		return nil
	case *ast.SwitchStmt:
		v.walk(1, n.Body)
		return nil
	case *ast.TypeSwitchStmt:
		v.walk(1, n.Body)
		return nil
	case *ast.FuncLit:
		v.walk(0, n.Body) // do not increment the complexity, just do the nesting
		return nil
	case *ast.BinaryExpr:
		v.complexity += v.binExpComplexity(n)
		return nil // skip visiting binexp sub-tree (already visited by binExpComplexity)
	case *ast.BranchStmt:
		if n.Label != nil {
			v.complexity += 1
		}
	}
	// TODO handle (at least) direct recursion

	return v
}

func (v *cognitiveComplexityVisitor) walk(complexityIncrement int, targets ...ast.Node) {
	v.complexity += complexityIncrement + v.nestingLevel
	nesting := v.nestingLevel
	v.nestingLevel++

	for _, t := range targets {
		if t == nil {
			continue
		}

		ast.Walk(v, t)
	}

	v.nestingLevel = nesting
}

func (cognitiveComplexityVisitor) binExpComplexity(n *ast.BinaryExpr) int {
	calculator := binExprComplexityCalculator{complexity: 0}
	ast.Walk(&calculator, n)

	return calculator.complexity
}

type binExprComplexityCalculator struct {
	complexity int
	currentOp  token.Token
}

func (v *binExprComplexityCalculator) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.BinaryExpr:
		isLogicOp := n.Op == token.LAND || n.Op == token.LOR
		if isLogicOp && n.Op != v.currentOp {
			v.complexity++
			v.currentOp = n.Op
		}
	case *ast.ParenExpr:
		v.complexity++
	}

	return v
}
