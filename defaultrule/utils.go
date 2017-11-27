package defaultrule

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
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

var commonMethods = map[string]bool{
	"Error":     true,
	"Read":      true,
	"ServeHTTP": true,
	"String":    true,
	"Write":     true,
}

func receiverType(fn *ast.FuncDecl) string {
	switch e := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	// The parser accepts much more than just the legal forms.
	return "invalid-type"
}

var knownNameExceptions = map[string]bool{
	"LastInsertId": true, // must match database/sql
	"kWh":          true,
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

var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

var allCapsRE = regexp.MustCompile(`^[A-Z0-9_]+$`)

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

var zeroLiteral = map[string]bool{
	"false": true, // bool
	// runes
	`'\x00'`: true,
	`'\000'`: true,
	// strings
	`""`: true,
	"``": true,
	// numerics
	"0":   true,
	"0.":  true,
	"0.0": true,
	"0i":  true,
}

func validType(T types.Type) bool {
	return T != nil &&
		T != types.Typ[types.Invalid] &&
		!strings.Contains(T.String(), "invalid type") // good but not foolproof
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}
