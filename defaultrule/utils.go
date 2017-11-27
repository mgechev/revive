package defaultrule

import (
	"strings"

	"github.com/mgechev/revive/file"
)

func isTest(f *file.File) bool {
	return strings.HasSuffix(f.Name, "_test.go")
}

const styleGuideBase = "https://golang.org/wiki/CodeReviewComments"
