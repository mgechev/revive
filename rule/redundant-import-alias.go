package rule

import (
	"fmt"
	"go/ast"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

const (
	defaultIgnoreUsed = false
)

// RedundantImportAlias lints given else constructs.
type RedundantImportAlias struct {
	ignoreUsed bool
	sync.Mutex
}

// Apply applies the rule to given file.
func (r *RedundantImportAlias) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configure(arguments)

	var failures []lint.Failure

	redundants := r.checkRedundantAliases(file.AST)

	for _, imp := range file.AST.Imports {
		if imp.Name == nil {
			continue
		}

		if _, exists := redundants[imp.Name.Name]; exists {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("Import alias \"%s\" is redundant", imp.Name.Name),
				Node:       imp,
				Category:   "imports",
			})
		}
	}
	return failures
}

// Name returns the rule name.
func (*RedundantImportAlias) Name() string {
	return "redundant-import-alias"
}

func (r *RedundantImportAlias) checkRedundantAliases(node ast.Node) map[string]string {

	var aliasedPackages = make(map[string]string)

	// First pass: Identify all aliases and their usage - by default alias is redundant
	ast.Inspect(node, func(n ast.Node) bool {
		imp, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}

		if imp.Name != nil && imp.Path != nil && imp.Name.Name != "_" && getImportPackageName(imp) == imp.Name.Name {
			aliasedPackages[imp.Name.Name] = "redundant"	
		}

		return true
	})

	if !r.ignoreUsed {
		return aliasedPackages
	}

	// Second pass: remove one-time used aliases
	ast.Inspect(node, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		x, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		// This alias is being used; it's not redundant
		if _, exists := aliasedPackages[x.Name]; exists {
			delete(aliasedPackages, x.Name)
		}

		return true
	})

	return aliasedPackages

}

func (r *RedundantImportAlias) configure(arguments lint.Arguments) {
	r.Lock()
	defer r.Unlock()

	r.ignoreUsed = defaultIgnoreUsed

	if len(arguments) == 0 {
		return
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		panic(fmt.Sprintf("Invalid argument to the redundant-import-alias rule. Expecting a k,v map, got %T", arguments[0]))
	}

	for k, v := range args {
		switch k {
		case "ignoreUsed":
			value, ok := v.(bool)
			if !ok {
				panic(fmt.Sprintf("Invalid argument to the redundant-import-alias rule, expecting string representation of an bool. Got '%v' (%T)", v, v))
			}
			r.ignoreUsed = value
		}
	}

}

func getImportPackageName(imp *ast.ImportSpec) string {
	path := strings.Trim(imp.Path.Value, `"`)
	parts := strings.Split(path, "/")
	packageName := parts[len(parts)-1]
	return packageName
}
