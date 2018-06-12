package rule

import (
	"fmt"
	"go/ast"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

type pkgMethods struct {
	lp *lint.Package
	mn map[string]map[string]bool
	mu *sync.Mutex
}

type packages struct {
	pkgs []pkgMethods
	mu   sync.Mutex
}

func (ps *packages) methodNames(lp *lint.Package) (map[string]map[string]bool, *sync.Mutex) {
	ps.mu.Lock()

	for _, pkg := range ps.pkgs {
		if pkg.lp == lp {
			ps.mu.Unlock()
			return pkg.mn, pkg.mu
		}
	}

	mn := make(map[string]map[string]bool)
	mu := sync.Mutex{}
	ps.pkgs = append(ps.pkgs, pkgMethods{lp: lp, mn: mn, mu: &mu})

	ps.mu.Unlock()
	return mn, &mu
}

var allPkgs = packages{pkgs: make([]pkgMethods, 1)}

// ConfusingNamingRule lints method names that differ only by capitalization
type ConfusingNamingRule struct{}

// Apply applies the rule to given file.
func (r *ConfusingNamingRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	mn, mu := allPkgs.methodNames(file.Pkg)
	walker := lintConfusingNames{
		methodNames: mn,
		mu:          mu,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(&walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ConfusingNamingRule) Name() string {
	return "confusing-naming"
}

//checkMethodName checks if a given method/function name is similar (just case differences) to other method/function of the same struct/file.
func checkMethodName(holder string, id *ast.Ident, w *lintConfusingNames) {
	name := strings.ToUpper(id.Name)
	w.mu.Lock()
	defer w.mu.Unlock()
	mnOfHolder := w.methodNames[holder]

	if mnOfHolder != nil {

		alreadyKnown := w.methodNames[holder][name]

		if alreadyKnown {
			// confusing names
			var msg string
			if holder == defaultStructName {
				msg = fmt.Sprintf("Function '%s' differs only by capitalization to other function in the same package", id.Name)
			} else {
				msg = fmt.Sprintf("Method '%s' differs only by capitalization to other method of '%s'", id.Name, holder)
			}
			w.onFailure(lint.Failure{
				Failure:    msg,
				Confidence: 1,
				Node:       id,
				Category:   "naming",
				URL:        "#TODO",
			})

			return
		}
	} else {
		w.methodNames[holder] = make(map[string]bool, 1)
	}

	// update the black list
	if w.methodNames[holder] == nil {
		println("no entry for '", holder, "'")
	}
	w.methodNames[holder][name] = true
}

type lintConfusingNames struct {
	mu          *sync.Mutex
	methodNames map[string]map[string]bool
	onFailure   func(lint.Failure)
}

const defaultStructName = "_" // used to map functions

//getStructName of a function receiver. Defaults to defaultStructName
func getStructName(r *ast.FieldList) string {
	result := defaultStructName

	if r == nil || len(r.List) < 1 {
		return result
	}

	t := r.List[0].Type

	if p, _ := t.(*ast.StarExpr); p != nil { // if a pointer receiver => dereference pointer receiver types
		t = p.X
	}

	if p, _ := t.(*ast.Ident); p != nil {
		result = p.Name
	}

	return result
}

func (w *lintConfusingNames) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.FuncDecl:
		// Exclude naming warnings for functions that are exported to C but
		// not exported in the Go API.
		// See https://github.com/golang/lint/issues/144.
		if ast.IsExported(v.Name.Name) || !isCgoExported(v) {
			checkMethodName(getStructName(v.Recv), v.Name, w)
		}
	default:
		// will add other checks like field names, struct names, etc.
	}

	return w
}
