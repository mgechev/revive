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
			expr, ok := goBuildExpr(comment.Text)
			if !ok {
				continue
			}

			// goVersion is the module's language version, e.g. "go1.21".
			goVersion := version.Lang("go" + file.Pkg.GoVersion().String())
			// Starting with Go 1.21 the go version in go.mod is a hard requirement, so a `//go:build go1.X` constraint
			// it already satisfies is redundant. Skip this check for Go 1.20 and earlier.
			if file.Pkg.IsAtLeastGoVersion(lint.Go121) && version.IsValid(expr) && version.Compare(goVersion, expr) >= 0 {
				return []lint.Failure{{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("The build tag %q is redundant for Go %s and can be removed", comment.Text, strings.TrimPrefix(goVersion, "go")),
				}}
			}

			if msg, ok := redundantFileNameConstraint(comment.Text, fileName, fileNameTags); ok {
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

const (
	goBuildPrefix   = "//go:build"
	plusBuildPrefix = "// +build"
)

// goBuildExpr returns the build expression of a `//go:build` comment line, e.g. "windows && amd64" for
// `//go:build windows && amd64`.
// The directive may be separated from the expression by any whitespace, so the prefix cannot be cut directly.
func goBuildExpr(line string) (string, bool) {
	if !constraint.IsGoBuild(line) {
		return "", false
	}

	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), goBuildPrefix)), true
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

// Name returns the rule name.
func (*RedundantBuildTagRule) Name() string {
	return "redundant-build-tag"
}

// redundantFileNameConstraint reports whether the `//go:build` line is wholly or partly implied by the GOOS/GOARCH
// suffix of the file name, together with the message describing the finding.
func redundantFileNameConstraint(line, fileName string, fileNameTags map[string]bool) (string, bool) {
	if len(fileNameTags) == 0 {
		return "", false
	}

	expr, err := constraint.Parse(line)
	if err != nil {
		return "", false // not a well-formed build constraint, nothing to say about it
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

// simplifyConstraint partially evaluates expr using the tags whose truth value is fixed by the file name.
// It returns a nil expression when expr is fully determined, in which case value is its result;
// otherwise it returns the simplified expression, which is expr itself when no simplification is possible.
func simplifyConstraint(expr constraint.Expr, tags map[string]bool) (simplified constraint.Expr, value bool) {
	switch e := expr.(type) {
	case *constraint.TagExpr:
		if value, ok := tags[e.Tag]; ok {
			return nil, value
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

// buildTagsFromFileName returns the build tags whose truth value is fixed by the GOOS/GOARCH suffix of fileName,
// e.g. "example_linux.go" fixes "linux" and "unix" to true and every other GOOS to false.
// It returns nil when the file name constrains nothing. Tags that remain undecided are absent from the result.
//
// The file name rules mirror goodOSArchFile:
// https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/go/build/build.go#1995
func buildTagsFromFileName(fileName string) map[string]bool {
	name, _, _ := strings.Cut(fileName, ".")

	i := strings.Index(name, "_")
	if i < 0 {
		return nil // everything before the first "_" is the file's own name, so "linux.go" constrains nothing
	}

	// Keep the leading "_" so that parts[0] holds the file's own name and is never a GOOS/GOARCH.
	parts := strings.Split(name[i:], "_")
	if n := len(parts); parts[n-1] == "test" {
		parts = parts[:n-1]
	}

	var goos, goarch string
	last := parts[len(parts)-1]
	switch n := len(parts); {
	case n >= 2 && knownOS[parts[n-2]] && knownArch[last]:
		goos, goarch = parts[n-2], last
	case knownOS[last]:
		goos = last
	case knownArch[last]:
		goarch = last
	default:
		return nil
	}

	tags := map[string]bool{}
	if goos != "" {
		for os := range knownOS {
			switch {
			case os == goos || slices.Contains(impliedOS[goos], os):
				tags[os] = true
			case slices.Contains(impliedOS[os], goos):
				// goos is a superset of os, e.g. a file named "*_linux.go" is also built for GOOS=android, so the
				// "android" tag stays undecided.
			default:
				tags[os] = false
			}
		}
		tags["unix"] = unixGOOS[goos]
	}
	if goarch != "" {
		for arch := range knownArch {
			tags[arch] = arch == goarch
		}
	}

	return tags
}

// impliedOS maps a GOOS value to the other GOOS build tags it also satisfies.
//
// Mirrors https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/go/build/build.go#1959
var impliedOS = map[string][]string{
	"android": {"linux"},
	"illumos": {"solaris"},
	"ios":     {"darwin"},
}

var (
	// knownOS is the list of past, present, and future known GOOS values.
	// Mirrors https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/internal/syslist/syslist.go#17
	knownOS = map[string]bool{
		"aix": true, "android": true, "darwin": true, "dragonfly": true, "freebsd": true,
		"hurd": true, "illumos": true, "ios": true, "js": true, "linux": true, "nacl": true,
		"netbsd": true, "openbsd": true, "plan9": true, "solaris": true, "wasip1": true,
		"windows": true, "zos": true,
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

	// unixGOOS is the set of GOOS values matched by the "unix" build tag.
	// Mirrors https://go.googlesource.com/go/+/refs/tags/go1.26.0/src/internal/syslist/syslist.go#41
	unixGOOS = map[string]bool{
		"aix":       true,
		"android":   true,
		"darwin":    true,
		"dragonfly": true,
		"freebsd":   true,
		"hurd":      true,
		"illumos":   true,
		"ios":       true,
		"linux":     true,
		"netbsd":    true,
		"openbsd":   true,
		"solaris":   true,
	}
)
