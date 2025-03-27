package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"

	"github.com/mgechev/revive/lint"
)

// TimeEqualRule shows where "==" and "!=" used for equality check time.Time
type TimeEqualRule struct {
	logger *slog.Logger
}

// Apply applies the rule to given file.
func (r *TimeEqualRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeEqual{file, onFailure}
	if err := w.file.Pkg.TypeCheck(); err != nil {
		r.logger.Info("TypeCheck returns error", "err", err)
		return nil
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*TimeEqualRule) Name() string {
	return "time-equal"
}

// SetLogger sets the logger field.
func (r *TimeEqualRule) SetLogger(logger *slog.Logger) {
	r.logger = logger.With("rule", r.Name())
}

type lintTimeEqual struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (l *lintTimeEqual) Visit(node ast.Node) ast.Visitor {
	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return l
	}

	switch expr.Op {
	case token.EQL, token.NEQ:
	default:
		return l
	}

	typeOfX := l.file.Pkg.TypeOf(expr.X)
	typeOfY := l.file.Pkg.TypeOf(expr.Y)
	bothAreOfTimeType := isNamedType(typeOfX, "time", "Time") && isNamedType(typeOfY, "time", "Time")
	if !bothAreOfTimeType {
		return l
	}

	negateStr := ""
	if token.NEQ == expr.Op {
		negateStr = "!"
	}

	l.onFailure(lint.Failure{
		Category:   lint.FailureCategoryTime,
		Confidence: 1,
		Node:       node,
		Failure:    fmt.Sprintf("use %s%s.Equal(%s) instead of %q operator", negateStr, gofmt(expr.X), gofmt(expr.Y), expr.Op),
	})

	return l
}
