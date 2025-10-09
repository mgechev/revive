package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/mgechev/revive/lint"
)

// exitFunctions is a map of std packages and functions that are considered as exit functions.
var exitFunctions = map[string]map[string]bool{
	"os":      {"Exit": true},
	"syscall": {"Exit": true},
	"log": {
		"Fatal":   true,
		"Fatalf":  true,
		"Fatalln": true,
		"Panic":   true,
		"Panicf":  true,
		"Panicln": true,
	},
	"flag": {"Parse": true},
}

func srcLine(src []byte, p token.Position) string {
	// Run to end of line in both directions if not at line start/end.
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}

// checkNumberOfArguments fails if the given number of arguments is not, at least, the expected one.
func checkNumberOfArguments(expected int, args lint.Arguments, ruleName string) error {
	if len(args) < expected {
		return fmt.Errorf("not enough arguments for %s rule, expected %d, got %d. Please check the rule's documentation", ruleName, expected, len(args))
	}
	return nil
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name from the argument. It is lowercased and without hyphens.
//
// Example: normalizeRuleOption("allowTypesBefore"), normalizeRuleOption("allow-types-before") -> "allowtypesbefore".
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

var normalizePathReplacer = strings.NewReplacer("-", "", "_", "", ".", "")

// normalizePath removes hyphens, underscores, and dots from the name
//
// Example: normalizePath("foo.bar-_buz") -> "foobarbuz".
func normalizePath(name string) string {
	return normalizePathReplacer.Replace(name)
}

// isVersionPath checks if a directory name is a version directory (v1, V2, etc.)
func isVersionPath(name string) bool {
	if len(name) < 2 || (name[0] != 'v' && name[0] != 'V') {
		return false
	}

	for i := 1; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}

	return true
}

var directiveCommentRE = regexp.MustCompile("^//(line |extern |export |[a-z0-9]+:[a-z0-9])") // see https://go-review.googlesource.com/c/website/+/442516/1..2/_content/doc/comment.md#494

func isDirectiveComment(line string) bool {
	return directiveCommentRE.MatchString(line)
}

// isCallToExitFunction checks if the function call is a call to an exit function.
func isCallToExitFunction(pkgName, functionName string) bool {
	return exitFunctions[pkgName] != nil && exitFunctions[pkgName][functionName]
}

var conditionalExitFunctions = map[string]map[string]func(*ast.CallExpr) bool{
    "flag": {
        "NewFlagSet": func(ce *ast.CallExpr) bool {
            if len(ce.Args) == 2 {
                if arg, ok := ce.Args[1].(*ast.SelectorExpr); ok {
                    if id, ok := arg.X.(*ast.Ident); ok {
                        return id.Name == "flag" && arg.Sel.Name == "ExitOnError"
                    }
                }
            }
            return false
        },
    },
}

// isConditionalExitFunction checks if the function call is a call to a conditional exit function.
func isConditionalExitFunction(pkgName, functionName string, ce *ast.CallExpr) bool {
	if m, ok := conditionalExitFunctions[pkgName]; ok {
		if check, ok := m[functionName]; ok {
			return check(ce)
		}
	}
	return false
}

// newInternalFailureError returns a slice of Failure with a single internal failure in it.
func newInternalFailureError(e error) []lint.Failure {
	return []lint.Failure{lint.NewInternalFailure(e.Error())}
}
