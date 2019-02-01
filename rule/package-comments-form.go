package rule

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/mgechev/revive/lint"
)

// PackageCommentsFormRule lints the package comments. It complains if
// the package comment is not of the right form.
type PackageCommentsFormRule struct{}

// Apply applies the rule to given file.
func (r *PackageCommentsFormRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if isTest(file) {
		return failures
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	fileAst := file.AST
	w := &lintPackageCommentsForm{fileAst, file, onFailure}
	ast.Walk(w, fileAst)
	return failures
}

// Name returns the rule name.
func (r *PackageCommentsFormRule) Name() string {
	return "package-comments-form"
}

type lintPackageCommentsForm struct {
	fileAst   *ast.File
	file      *lint.File
	onFailure func(lint.Failure)
}

func (l *lintPackageCommentsForm) Visit(_ ast.Node) ast.Visitor {
	if l.file.IsTest() {
		return nil
	}

	const ref = styleGuideBase + "#package-comments"
	prefix := "Package " + l.fileAst.Name.Name + " "

	if l.fileAst.Doc != nil {
		s := l.fileAst.Doc.Text()
		if ts := strings.TrimLeft(s, " \t"); ts != s {
			l.onFailure(lint.Failure{
				Category:   "comments",
				Node:       l.fileAst.Doc,
				Confidence: 1,
				Failure:    "package comment should not have leading space",
			})
			s = ts
		}
		// Only non-main packages need to keep to this form.
		if !l.file.Pkg.IsMain() && !strings.HasPrefix(s, prefix) {
			l.onFailure(lint.Failure{
				Category:   "comments",
				Node:       l.fileAst.Doc,
				Confidence: 1,
				Failure:    fmt.Sprintf(`package comment should be of the form "%s..."`, prefix),
			})
		}
	}
	return nil
}
