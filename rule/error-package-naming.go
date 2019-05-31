package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
	"strings"
	"regexp"
)

// ErrorPackageNamingRule lints package name
type ErrorPackageNamingRule struct{}

// Apply applies the rule to given file.
func (r *ErrorPackageNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if isTest(file) {
		return failures
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	fileAst := file.AST
	w := &lintErrorPackageNaming{fileAst, file, onFailure}

	ast.Walk(w, fileAst)
	return failures
}

// Name returns the rule name.
func (r *ErrorPackageNamingRule) Name() string {
	return "error-package-naming"
}

type lintErrorPackageNaming struct {
	fileAst   *ast.File
	file      *lint.File
	onFailure func(lint.Failure)
}

func (l *lintErrorPackageNaming) Visit(n ast.Node) ast.Visitor {
	if l.file.IsTest() {
		return nil
	}

	pkgName := l.fileAst.Name.Name

	if strings.Contains(pkgName, "_") && !strings.HasSuffix(pkgName, "_test") {
		l.onFailure(lint.Failure{
			Category:   "naming",
			Node:       l.fileAst.Doc,
			Confidence: 1,
			Failure:    fmt.Sprintf(`don't use an underscore in package name`,),
		})
	}

	anyCapsRE := regexp.MustCompile(`[A-Z]`)
	if anyCapsRE.MatchString(pkgName) {
		l.onFailure(lint.Failure{
			Category:  "mixed-caps",
			Node:       l.fileAst.Doc,
			Confidence: 1,
			Failure:    fmt.Sprintf("don't use MixedCaps in package name; %s should be %s\n", pkgName, strings.ToLower(pkgName)),
		})
	}
	return nil
}

