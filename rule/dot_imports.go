package rule

import (
	"fmt"
	"go/ast"
	"sync"

	"github.com/mgechev/revive/lint"
)

// DotImportsRule forbids . imports.
type DotImportsRule struct {
	allowedPackages allowPackages

	configureOnce sync.Once
}

// Apply applies the rule to given file.
func (r *DotImportsRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var configureErr error
	r.configureOnce.Do(func() { configureErr = r.configure(arguments) })

	if configureErr != nil {
		return []lint.Failure{lint.NewInternalFailure(configureErr.Error())}
	}

	var failures []lint.Failure

	fileAst := file.AST
	walker := lintImports{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
		allowPackages: r.allowedPackages,
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*DotImportsRule) Name() string {
	return "dot-imports"
}

func (r *DotImportsRule) configure(arguments lint.Arguments) error {
	r.allowedPackages = allowPackages{}
	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the dot-imports rule. Expecting a k,v map, got %T", arguments[0])
	}

	if allowedPkgArg, ok := args["allowedPackages"]; ok {
		pkgs, ok := allowedPkgArg.([]any)
		if !ok {
			return fmt.Errorf("invalid argument to the dot-imports rule, []string expected. Got '%v' (%T)", allowedPkgArg, allowedPkgArg)
		}
		for _, p := range pkgs {
			pkg, ok := p.(string)
			if !ok {
				return fmt.Errorf("invalid argument to the dot-imports rule, string expected. Got '%v' (%T)", p, p)
			}
			r.allowedPackages.add(pkg)
		}
	}
	return nil
}

type lintImports struct {
	file          *lint.File
	fileAst       *ast.File
	onFailure     func(lint.Failure)
	allowPackages allowPackages
}

func (w lintImports) Visit(_ ast.Node) ast.Visitor {
	for _, importSpec := range w.fileAst.Imports {
		isDotImport := importSpec.Name != nil && importSpec.Name.Name == "."
		if isDotImport && !w.allowPackages.isAllowedPackage(importSpec.Path.Value) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Failure:    "should not use dot imports",
				Node:       importSpec,
				Category:   "imports",
			})
		}
	}
	return nil
}

type allowPackages map[string]struct{}

func (ap allowPackages) add(pkg string) {
	ap[fmt.Sprintf(`"%s"`, pkg)] = struct{}{} // import path strings are with double quotes
}

func (ap allowPackages) isAllowedPackage(pkg string) bool {
	_, allowed := ap[pkg]
	return allowed
}
