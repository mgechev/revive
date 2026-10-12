package rule

import (
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/version"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mgechev/revive/lint"
)

// RedundantBuildTagRule lints the presence of redundant build tags.
type RedundantBuildTagRule struct{}

// Apply triggers on three kinds of redundant build tags:
//   - a `// +build` tag found after a `//go:build` tag, since `//go:build` comments are automatically added by gofmt
//     when Go 1.17+ is used;
//   - a `//go:build go1.X` version constraint that is already guaranteed by the go version in go.mod
//     (Go 1.21+, where the go directive is a hard requirement);
//   - a `//go:build` constraint that is already implied by the GOOS/GOARCH suffix of the file name,
//     e.g. `//go:build windows` in a file named `example_windows.go`.
//
// See https://pkg.go.dev/cmd/go#hdr-Build_constraints
func (*RedundantBuildTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	fileName := filepath.Base(file.Name)
	fileNameTags := buildTagsFromFileName(fileName)
	// goVersion is the module's language version, e.g. "go1.21".
	goVersion := version.Lang("go" + file.Pkg.GoVersion().String())
	// Starting with Go 1.21 the go version in go.mod is a hard requirement, so a `//go:build go1.X` constraint
	// it already satisfies is redundant. Skip this check for Go 1.20 and earlier.
	checkGoVersion := file.Pkg.IsAtLeastGoVersion(lint.Go121)

	for _, group := range file.AST.Comments {
		if group.End() > file.AST.Package {
			break // build constraints must appear before the package clause
		}

		// The `// +build` tag is reported first: when a file carries both tags it is the only actionable finding,
		// since dropping a `//go:build` line that `// +build` still backs up just makes gofmt add it again.
		if comment, ok := redundantPlusBuildTag(group); ok {
			return []lint.Failure{{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       comment,
				Failure:    fmt.Sprintf("The build tag %q is redundant since Go 1.17 and can be removed", plusBuildPrefix),
			}}
		}

		for _, comment := range group.List {
			if !constraint.IsGoBuild(comment.Text) {
				continue
			}

			expr, err := constraint.Parse(comment.Text)
			if err != nil {
				continue // not a well-formed build constraint, nothing to say about it
			}

			if checkGoVersion && satisfiedGoVersionTag(expr, goVersion) {
				return []lint.Failure{{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("The build tag %q is redundant for Go %s and can be removed", comment.Text, strings.TrimPrefix(goVersion, "go")),
				}}
			}

			if msg, ok := redundantFileNameConstraint(expr, comment.Text, fileName, fileNameTags); ok {
				return []lint.Failure{{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    msg,
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

const (
	goBuildPrefix   = "//go:build"
	plusBuildPrefix = "// +build"
)

// satisfiedGoVersionTag reports whether expr is a lone `go1.X` version tag that goVersion already satisfies.
func satisfiedGoVersionTag(expr constraint.Expr, goVersion string) bool {
	tag, ok := expr.(*constraint.TagExpr)

	return ok && version.IsValid(tag.Tag) && version.Compare(goVersion, tag.Tag) >= 0
}

// redundantPlusBuildTag returns the `// +build` comment of the group when it follows a `//go:build` one,
// which makes it redundant since Go 1.17.
func redundantPlusBuildTag(group *ast.CommentGroup) (*ast.Comment, bool) {
	hasGoBuild := false
	for _, comment := range group.List {
		switch {
		case constraint.IsGoBuild(comment.Text):
			hasGoBuild = true
		case hasGoBuild && constraint.IsPlusBuild(comment.Text):
			return comment, true
		}
	}

	return nil, false
}

// redundantFileNameConstraint reports whether the `//go:build` line is wholly or partly implied by the GOOS/GOARCH
// suffix of the file name, together with the message describing the finding.
func redundantFileNameConstraint(expr constraint.Expr, line, fileName string, fileNameTags []string) (string, bool) {
	if len(fileNameTags) == 0 {
		return "", false
	}

	switch simplified, value := simplifyConstraint(expr, fileNameTags); {
	case simplified == nil && value:
		return fmt.Sprintf("The build tag %q is implied by the file name %q and can be removed", line, fileName), true
	case simplified == nil && !value:
		return fmt.Sprintf("The build tag %q contradicts the file name %q, so the file is never built", line, fileName), true
	case simplified != expr:
		return fmt.Sprintf("The build tag %q is partly implied by the file name %q and can be reduced to %q", line, fileName, goBuildPrefix+" "+simplified.String()), true
	}

	return "", false
}

// simplifyConstraint partially evaluates expr using the tags that the file name guarantees to be satisfied.
// It returns a nil expression when expr is fully determined, in which case value is its result;
// otherwise it returns the simplified expression, which is expr itself when no simplification is possible.
func simplifyConstraint(expr constraint.Expr, tags []string) (simplified constraint.Expr, value bool) {
	switch e := expr.(type) {
	case *constraint.TagExpr:
		if slices.Contains(tags, e.Tag) {
			return nil, true
		}

	case *constraint.NotExpr:
		x, xValue := simplifyConstraint(e.X, tags)
		switch {
		case x == nil:
			return nil, !xValue
		case x != e.X:
			return &constraint.NotExpr{X: x}, false
		}

	case *constraint.AndExpr:
		x, xValue := simplifyConstraint(e.X, tags)
		y, yValue := simplifyConstraint(e.Y, tags)
		switch {
		case x == nil && !xValue, y == nil && !yValue:
			return nil, false
		case x == nil && y == nil:
			return nil, true
		case x == nil:
			return y, false
		case y == nil:
			return x, false
		case x != e.X || y != e.Y:
			return &constraint.AndExpr{X: x, Y: y}, false
		}

	case *constraint.OrExpr:
		x, xValue := simplifyConstraint(e.X, tags)
		y, yValue := simplifyConstraint(e.Y, tags)
		switch {
		case x == nil && xValue, y == nil && yValue:
			return nil, true
		case x == nil && y == nil:
			return nil, false
		case x == nil:
			return y, false
		case y == nil:
			return x, false
		case x != e.X || y != e.Y:
			return &constraint.OrExpr{X: x, Y: y}, false
		}
	}

	return expr, false
}

// buildTagsFromFileName returns the build tags that the GOOS/GOARCH suffix of fileName guarantees to be satisfied,
// e.g. "example_linux_amd64.go" guarantees "linux" and "amd64". It returns nil when the file name constrains nothing.
//
// Only the suffix tags themselves are guaranteed: goodOSArchFile matches them with the same matcher as any other
// build tag, so e.g. `-tags=linux` builds "example_linux.go" for every GOOS. Nothing can therefore be concluded
// about the "unix" tag, about a GOOS implying another one (`GOOS=android` satisfies "linux"), or about the
// remaining GOOS/GOARCH tags being unsatisfied.
//
// Mirrors goodOSArchFile: https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/go/build/build.go#1995
func buildTagsFromFileName(fileName string) []string {
	name, _, _ := strings.Cut(fileName, ".")

	i := strings.Index(name, "_")
	if i < 0 {
		return nil // everything before the first "_" is the file's own name, so "linux.go" constrains nothing
	}

	// Keep the leading "_" so that parts[0] is the empty string and can never be a GOOS/GOARCH.
	parts := strings.Split(name[i:], "_")
	if n := len(parts); parts[n-1] == "test" {
		parts = parts[:n-1]
	}

	last := parts[len(parts)-1]
	switch n := len(parts); {
	case n >= 2 && knownOS[parts[n-2]] && knownArch[last]:
		return []string{parts[n-2], last}
	case knownOS[last], knownArch[last]:
		return []string{last}
	}

	return nil
}

var (
	// knownOS is the list of past, present, and future known GOOS values.
	// Mirrors https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/internal/syslist/syslist.go#17
	knownOS = map[string]bool{
		"aix":       true,
		"android":   true,
		"darwin":    true,
		"dragonfly": true,
		"freebsd":   true,
		"hurd":      true,
		"illumos":   true,
		"ios":       true,
		"js":        true,
		"linux":     true,
		"nacl":      true,
		"netbsd":    true,
		"openbsd":   true,
		"plan9":     true,
		"solaris":   true,
		"wasip1":    true,
		"windows":   true,
		"zos":       true,
	}

	// knownArch is the list of past, present, and future known GOARCH values.
	// Mirrors https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/internal/syslist/syslist.go#58
	knownArch = map[string]bool{
		"386":         true,
		"amd64":       true,
		"amd64p32":    true,
		"arm":         true,
		"armbe":       true,
		"arm64":       true,
		"arm64be":     true,
		"loong64":     true,
		"mips":        true,
		"mipsle":      true,
		"mips64":      true,
		"mips64le":    true,
		"mips64p32":   true,
		"mips64p32le": true,
		"ppc":         true,
		"ppc64":       true,
		"ppc64le":     true,
		"riscv":       true,
		"riscv64":     true,
		"s390":        true,
		"s390x":       true,
		"sparc":       true,
		"sparc64":     true,
		"wasm":        true,
	}
)
