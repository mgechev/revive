package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"github.com/mgechev/revive/lint"
	"golang.org/x/tools/go/ast/astutil"
)

// CognitiveComplexityRule lints given else constructs.
type CognitiveComplexityRule struct {
	maxComplexity int
	sync.Mutex
}

const defaultMaxCognitiveComplexity = 7

func (r *CognitiveComplexityRule) configure(arguments lint.Arguments) {
	r.Lock()
	defer r.Unlock()
	if r.maxComplexity == 0 {

		if len(arguments) < 1 {
			r.maxComplexity = defaultMaxCognitiveComplexity
			return
		}

		complexity, ok := arguments[0].(int64)
		if !ok {
			panic(fmt.Sprintf("invalid argument type for cognitive-complexity, expected int64, got %T", arguments[0]))
		}
		r.maxComplexity = int(complexity)
	}
}

// Apply applies the rule to given file.
func (r *CognitiveComplexityRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configure(arguments)

	var failures []lint.Failure

	linter := cognitiveComplexityLinter{
		file:          file,
		maxComplexity: r.maxComplexity,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	linter.lintCognitiveComplexity()

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

func (r cognitiveComplexityLinter) lintCognitiveComplexity() {
	f := r.file
	for _, decl := range f.AST.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
			v := cognitiveComplexityVisitor{}
			c := v.subTreeComplexity(fn.Body)
			if c > r.maxComplexity {
				r.onFailure(lint.Failure{
					Confidence: 1,
					Category:   "maintenance",
					Failure:    fmt.Sprintf("function %s has cognitive complexity %d (> max enabled %d)", funcName(fn), c, r.maxComplexity),
					Node:       fn,
				})
			}
		}
	}
}

type cognitiveComplexityVisitor struct {
	complexity   int
	nestingLevel int
}

// subTreeComplexity calculates the cognitive complexity of an AST-subtree.
func (r cognitiveComplexityVisitor) subTreeComplexity(n ast.Node) int {
	ast.Walk(&r, n)
	return r.complexity
}

// Visit implements the ast.Visitor interface.
func (r *cognitiveComplexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.IfStmt:
		targets := []ast.Node{n.Cond, n.Body, n.Else}
		r.walk(1, targets...)
		return nil
	case *ast.ForStmt:
		targets := []ast.Node{n.Cond, n.Body}
		r.walk(1, targets...)
		return nil
	case *ast.RangeStmt:
		r.walk(1, n.Body)
		return nil
	case *ast.SelectStmt:
		r.walk(1, n.Body)
		return nil
	case *ast.SwitchStmt:
		r.walk(1, n.Body)
		return nil
	case *ast.TypeSwitchStmt:
		r.walk(1, n.Body)
		return nil
	case *ast.FuncLit:
		r.walk(0, n.Body) // do not increment the complexity, just do the nesting
		return nil
	case *ast.BinaryExpr:
		r.complexity += r.binExpComplexity(n)
		return nil // skip visiting binexp sub-tree (already visited by binExpComplexity)
	case *ast.BranchStmt:
		if n.Label != nil {
			r.complexity++
		}
	}
	// TODO handle (at least) direct recursion

	return r
}

func (r *cognitiveComplexityVisitor) walk(complexityIncrement int, targets ...ast.Node) {
	r.complexity += complexityIncrement + r.nestingLevel
	nesting := r.nestingLevel
	r.nestingLevel++

	for _, t := range targets {
		if t == nil {
			continue
		}

		ast.Walk(r, t)
	}

	r.nestingLevel = nesting
}

func (r cognitiveComplexityVisitor) binExpComplexity(n *ast.BinaryExpr) int {
	calculator := binExprComplexityCalculator{opsStack: []token.Token{}}

	astutil.Apply(n, calculator.pre, calculator.post)

	return calculator.complexity
}

type binExprComplexityCalculator struct {
	complexity    int
	opsStack      []token.Token // stack of bool operators
	subexpStarted bool
}

func (r *binExprComplexityCalculator) pre(c *astutil.Cursor) bool {
	switch n := c.Node().(type) {
	case *ast.BinaryExpr:
		isBoolOp := n.Op == token.LAND || n.Op == token.LOR
		if !isBoolOp {
			break
		}

		ops := len(r.opsStack)
		// if
		// 		is the first boolop in the expression OR
		// 		is the first boolop inside a subexpression (...) OR
		//		is not the same to the previous one
		// then
		//      increment complexity
		if ops == 0 || r.subexpStarted || n.Op != r.opsStack[ops-1] {
			r.complexity++
			r.subexpStarted = false
		}

		r.opsStack = append(r.opsStack, n.Op)
	case *ast.ParenExpr:
		r.subexpStarted = true
	}

	return true
}

func (r *binExprComplexityCalculator) post(c *astutil.Cursor) bool {
	switch n := c.Node().(type) {
	case *ast.BinaryExpr:
		isBoolOp := n.Op == token.LAND || n.Op == token.LOR
		if !isBoolOp {
			break
		}

		ops := len(r.opsStack)
		if ops > 0 {
			r.opsStack = r.opsStack[:ops-1]
		}
	case *ast.ParenExpr:
		r.subexpStarted = false
	}

	return true
}
