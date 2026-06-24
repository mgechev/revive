package rule

import (
	"fmt"
	"go/version"
	"strings"

	"github.com/mgechev/revive/lint"
)

// RedundantBuildTagRule lints the presence of redundant build tags.
type RedundantBuildTagRule struct{}

// Apply triggers on two kinds of redundant build tags:
//   - an old `// +build` tag found after a new `//go:build` tag, since `//go:build`
//     comments are automatically added by gofmt when Go 1.17+ is used;
//   - a `//go:build go1.X` version constraint that is already guaranteed by the
//     go version in go.mod (Go 1.21+, where the go directive is a hard requirement).
//
// See https://pkg.go.dev/cmd/go#hdr-Build_constraints
func (*RedundantBuildTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	for _, group := range file.AST.Comments {
		hasGoBuild := false
		for _, comment := range group.List {
			if ver, ok := strings.CutPrefix(comment.Text, "//go:build "); ok {
				hasGoBuild = true

				// Starting with Go 1.21, the go version in go.mod is a hard requirement.
				// Ignore this check for Go 1.20 and earlier.
				if file.Pkg.IsAtLeastGoVersion(lint.Go121) && version.IsValid(ver) {
					segments := file.Pkg.GoVersion().Segments()
					fileVersion := fmt.Sprintf("%d.%d", segments[0], segments[1])
					if version.Compare("go"+fileVersion, ver) >= 0 {
						return []lint.Failure{{
							Category:   lint.FailureCategoryStyle,
							Confidence: 1,
							Node:       comment,
							Failure:    fmt.Sprintf("The build tag %q is redundant for Go %s and can be removed", comment.Text, fileVersion),
						}}
					}
				}

				continue
			}

			const oldGoBuildPrefix = "// +build"
			if hasGoBuild && strings.HasPrefix(comment.Text, oldGoBuildPrefix) {
				return []lint.Failure{{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("The build tag %q is redundant since Go 1.17 and can be removed", oldGoBuildPrefix),
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
