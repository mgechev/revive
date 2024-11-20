package rule

import (
	"strings"

	"github.com/mgechev/revive/lint"
)

// RedundantBuildTagRule lints the presence of redundant build tags.
type RedundantBuildTagRule struct{}

// Apply triggers if an old build tag `// +build` is found after a new one `//go:build`.
// `//go:build` comments are automatically added by gofmt when Go 1.17+ is used.
// See https://pkg.go.dev/cmd/go#hdr-Build_constraints
func (*RedundantBuildTagRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		hasGoBuild := false
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, "//go:build ") {
				hasGoBuild = true
				continue
			}
			if hasGoBuild && strings.HasPrefix(comment.Text, "// +build ") {
				failures = append(failures, lint.Failure{
					Category:   "style",
					Confidence: 1,
					Node:       comment,
					Failure:    `The build tag "// +build" is redundant since Go 1.17 and can be removed`,
				})
				return failures
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*RedundantBuildTagRule) Name() string {
	return "redundant-build-tag"
}
