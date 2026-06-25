package rule

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/mgechev/revive/lint"
)

// canonicalImportCommentRegexp matches a canonical import comment body, e.g.: import "rsc.io/pdf" or import `rsc.io/pdf`.
var canonicalImportCommentRegexp = regexp.MustCompile("^import\\s+(?:\"[^\"]+\"|`[^`]+`)$")

// RedundantCanonicalImport warns on canonical import path comments that are redundant in module mode (Go 1.11+).
// See https://go.dev/doc/go1.4#canonicalimports.
type RedundantCanonicalImport struct{}

// Apply applies the rule to given file.
func (*RedundantCanonicalImport) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go111) {
		return nil
	}

	packageLine := file.ToPosition(file.AST.Name.End()).Line

	for _, cg := range file.AST.Comments {
		for _, c := range cg.List {
			if isBeforePackageName(c, file) {
				continue
			}
			if file.ToPosition(c.Pos()).Line > packageLine {
				return nil // past the package clause line; comments are ordered by position
			}

			if !canonicalImportCommentRegexp.MatchString(commentBody(c.Text)) {
				continue
			}

			return []lint.Failure{
				{
					Category:   lint.FailureCategoryImports,
					Node:       c,
					Confidence: 1,
					Failure:    fmt.Sprintf("Canonical import comment `%s` is redundant and can be removed", c.Text),
				},
			}
		}
	}

	return nil
}

// isBeforePackageName reports whether the comment appears before the package name.
func isBeforePackageName(c *ast.Comment, file *lint.File) bool {
	return c.Pos() < file.AST.Name.End()
}

// commentBody strips the comment markers from c.Text, handling both line and block comments.
func commentBody(text string) string {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	return strings.TrimSpace(text)
}

// Name returns the rule name.
func (*RedundantCanonicalImport) Name() string {
	return "redundant-canonical-import"
}
