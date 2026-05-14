package rule

import (
	"fmt"
	"go/version"
	"strings"

	"github.com/mgechev/revive/lint"
)

// RedundantBuildTagRule lints the presence of redundant build tags.
type RedundantBuildTagRule struct{}

// Apply triggers if an old build tag `// +build` is found after a new one `//go:build`.
// `//go:build` comments are automatically added by gofmt when Go 1.17+ is used.
// See https://pkg.go.dev/cmd/go#hdr-Build_constraints
func (*RedundantBuildTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	for _, group := range file.AST.Comments {
		hasGoBuild := false
		for _, comment := range group.List {
			if ver, ok := strings.CutPrefix(comment.Text, "//go:build "); ok {
				hasGoBuild = true

				// Starting from the Go 1.21, the go version in go.mod is a hard requirement.
				// Ignore this check for Go 1.20 and earlier.
				if file.Pkg.IsAtLeastGoVersion(lint.Go121) && version.IsValid(ver) {
					fileVersion := file.Pkg.GoVersion().String()
					if version.Compare("go"+fileVersion, ver) >= 0 {
						return []lint.Failure{{
							Category:   lint.FailureCategoryStyle,
							Confidence: 1,
							Node:       comment,
							Failure:    fmt.Sprintf(`The build tag %q is redundant for Go %s and can be removed`, comment.Text, fileVersion),
						}}
					}
				}

				continue
			}

			if hasGoBuild && strings.HasPrefix(comment.Text, "// +build ") {
				return []lint.Failure{{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    `The build tag "// +build" is redundant since Go 1.17 and can be removed`,
				}}
			}
		}
	}

	return []lint.Failure{}
}

// Name returns the rule name.
func (*RedundantBuildTagRule) Name() string {
	return "redundant-build-tag"
}
