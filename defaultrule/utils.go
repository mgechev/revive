package defaultrule

import (
	"go/ast"
	"strings"

	"github.com/mgechev/revive/file"
)

const styleGuideBase = "https://golang.org/wiki/CodeReviewComments"

// isBlank returns whether id is the blank identifier "_".
// If id == nil, the answer is false.
func isBlank(id *ast.Ident) bool { return id != nil && id.Name == "_" }

func isTest(f *file.File) bool {
	return strings.HasSuffix(f.Name, "_test.go")
}
