package rule

import (
	"fmt"
	"go/ast"
	"log"
	"regexp"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ErrorfRule suggests using `fmt.Errorf` instead of `errors.New(fmt.Sprintf())`.
type ErrorfRule struct {
	logger *log.Logger
}

// Apply applies the rule to given file.
func (r *ErrorfRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	walker := lintErrorf{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	if err := file.Pkg.TypeCheck(); err != nil {
		r.logger.Printf("Rule=%q TypeCheck() error=%v\n", r.Name(), err)
	}
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*ErrorfRule) Name() string {
	return "errorf"
}

// SetLogger sets the logger field.
func (r *ErrorfRule) SetLogger(logger *log.Logger) {
	r.logger = logger
}

type lintErrorf struct {
	file      *lint.File
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

func (w lintErrorf) Visit(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if !ok || len(ce.Args) != 1 {
		return w
	}
	isErrorsNew := isPkgDot(ce.Fun, "errors", "New")
	var isTestingError bool
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if ok && se.Sel.Name == "Error" {
		if typ := w.file.Pkg.TypeOf(se.X); typ != nil {
			isTestingError = typ.String() == "*testing.T"
		}
	}
	if !isErrorsNew && !isTestingError {
		return w
	}
	arg := ce.Args[0]
	ce, ok = arg.(*ast.CallExpr)
	if !ok || !isPkgDot(ce.Fun, "fmt", "Sprintf") {
		return w
	}
	errorfPrefix := "fmt"
	if isTestingError {
		errorfPrefix = w.file.Render(se.X)
	}

	failure := lint.Failure{
		Category:   lint.FailureCategoryErrors,
		Node:       n,
		Confidence: 1,
		Failure:    fmt.Sprintf("should replace %s(fmt.Sprintf(...)) with %s.Errorf(...)", w.file.Render(se), errorfPrefix),
	}

	m := srcLineWithMatch(w.file, ce, `^(.*)`+w.file.Render(se)+`\(fmt\.Sprintf\((.*)\)\)(.*)$`)
	if m != nil {
		failure.ReplacementLine = m[1] + errorfPrefix + ".Errorf(" + m[2] + ")" + m[3]
	}

	w.onFailure(failure)

	return w
}

func srcLineWithMatch(file *lint.File, node ast.Node, pattern string) (m []string) {
	line := srcLine(file.Content(), file.ToPosition(node.Pos()))
	line = strings.TrimSuffix(line, "\n")
	rx := regexp.MustCompile(pattern)
	return rx.FindStringSubmatch(line)
}
