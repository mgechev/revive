package rule

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"regexp"

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
}

func isCgoExported(f *ast.FuncDecl) bool {
	if f.Recv != nil || f.Doc == nil {
		return false
	}

	cgoExport := regexp.MustCompile(fmt.Sprintf("(?m)^//export %s$", regexp.QuoteMeta(f.Name.Name)))
	for _, c := range f.Doc.List {
		if cgoExport.MatchString(c.Text) {
			return true
		}
	}
	return false
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// isPkgDot checks if the expression is <pkg>.<name>
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
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

// pick yields a list of nodes by picking them from a sub-ast with root node n.
// Nodes are selected by applying the fselect function
func pick(n ast.Node, fselect func(n ast.Node) bool) []ast.Node {
	var result []ast.Node

	if n == nil {
		return result
	}

	onSelect := func(n ast.Node) {
		result = append(result, n)
	}
	p := picker{fselect: fselect, onSelect: onSelect}
	ast.Walk(p, n)
	return result
}

type picker struct {
	fselect  func(n ast.Node) bool
	onSelect func(n ast.Node)
}

func (p picker) Visit(node ast.Node) ast.Visitor {
	if p.fselect == nil {
		return nil
	}

	if p.fselect(node) {
		p.onSelect(node)
	}

	return p
}

// gofmt returns a string representation of an AST subtree.
func gofmt(x any) string {
	buf := bytes.Buffer{}
	fs := token.NewFileSet()
	printer.Fprint(&buf, fs, x)
	return buf.String()
}

// checkNumberOfArguments fails if the given number of arguments is not, at least, the expected one
func checkNumberOfArguments(expected int, args lint.Arguments, ruleName string) error {
	if len(args) < expected {
		return fmt.Errorf("not enough arguments for %s rule, expected %d, got %d. Please check the rule's documentation", ruleName, expected, len(args))
	}
	return nil
}

var directiveCommentRE = regexp.MustCompile("^//(line |extern |export |[a-z0-9]+:[a-z0-9])") // see https://go-review.googlesource.com/c/website/+/442516/1..2/_content/doc/comment.md#494

func isDirectiveComment(line string) bool {
	return directiveCommentRE.MatchString(line)
}

// isCallToExitFunction checks if the function call is a call to an exit function.
func isCallToExitFunction(pkgName, functionName string) bool {
	return exitFunctions[pkgName] != nil && exitFunctions[pkgName][functionName]
}

// newInternalFailureError returns a slice of Failure with a single internal failure in it
func newInternalFailureError(e error) []lint.Failure {
	return []lint.Failure{lint.NewInternalFailure(e.Error())}
}

// anyCaps checks if string contains any capital letter
func anyCaps(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// allCaps checks if string have all capital letters
func allCaps(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

// isUpperCaseConst checks if string is in constant name format like `SOME_CONST`, `SOME_CONST_2`, `X123_3`, `_SOME_PRIVATE_CONST` (#851, #865)
func isUpperCaseConst(s string) bool {
	if len(s) == 0 {
		return false
	}

	i := 0

	// Optional leading underscore
	if s[0] == '_' {
		i++
		if i >= len(s) {
			return false // underscore only is not valid
		}
	}

	// Must start with an uppercase letter
	if s[i] < 'A' || s[i] > 'Z' {
		return false
	}
	i++

	// Continue parsing the rest
	for i < len(s) {
		c := s[i]
		if c == '_' {
			// Underscore must be followed by at least one uppercase letter or digit
			i++
			if i >= len(s) {
				return false
			}
			if !isUpperOrDigit(s[i]) {
				return false
			}
		} else if !isUpperOrDigit(c) {
			return false
		}
		i++
	}

	return true
}

// isUpperOrDigit checks if character is capital or digit
func isUpperOrDigit(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
