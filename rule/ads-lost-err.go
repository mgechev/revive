package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ADSLostErrRule lints program exit at functions other than main or init.
type ADSLostErrRule struct{}

// Apply applies the rule to given file.
func (r *ADSLostErrRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintLostErr{onFailure, "errors", "New"}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *ADSLostErrRule) Name() string {
	return "ads-newerr"
}

type lintLostErr struct {
	onFailure  func(lint.Failure)
	targetPkg  string
	targetFunc string
}

func (w lintLostErr) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	pkg, fn := getPkgFunc(ce)

	if pkg == w.targetPkg && fn == w.targetFunc {
		if pkg == "errors" && fn == "New" {
			w.targetFunc = "MessageOption"
			return w
		}

		s := searchErr{&w, ce}
		for _, exp := range ce.Args {
			ast.Walk(s, exp)
		}

		w.targetFunc = "New"
	}

	return w
}

type searchErr struct {
	w  *lintLostErr
	fc *ast.CallExpr
}

func (s searchErr) Visit(node ast.Node) ast.Visitor {

	id, ok := node.(*ast.Ident)
	if !ok {
		return s
	}
	if id.Name == "err" {
		s.w.onFailure(lint.Failure{
			Confidence: 0.8,
			Node:       s.fc,
			Category:   "bad practice",
			Failure:    "original error is lost, consider using errors.NewFromError",
		})

	}

	return s
}

func getPkgFunc(ce *ast.CallExpr) (string, string) {
	fc, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", ""
	}
	id, ok := fc.X.(*ast.Ident)
	if !ok {
		return "", ""
	}

	return id.Name, fc.Sel.Name
}
