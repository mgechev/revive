package rule

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/mgechev/revive/lint"
)

// RedundantImportAlias warns on import aliases matching the imported package name.
type RedundantImportAlias struct{}

// Apply applies the rule to given file.
func (*RedundantImportAlias) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	_ = file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()

	for _, imp := range file.AST.Imports {
		if imp.Name == nil {
			continue
		}

		if getImportPackageName(imp, typesInfo) == imp.Name.Name {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("Import alias %q is redundant", imp.Name.Name),
				Node:       imp,
				Category:   lint.FailureCategoryImports,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*RedundantImportAlias) Name() string {
	return "redundant-import-alias"
}

func getImportPackageName(imp *ast.ImportSpec, typesInfo *types.Info) string {
	if typesInfo != nil {
		if obj, ok := typesInfo.Defs[imp.Name]; ok {
			if pkgName, ok := obj.(*types.PkgName); ok {
				if imported := pkgName.Imported(); imported != nil && imported.Name() != "" {
					return imported.Name()
				}
			}
		}
	}

	const pathSep = "/"
	const strDelim = `"`

	path := imp.Path.Value
	i := strings.LastIndex(path, pathSep)
	if i == -1 {
		return strings.Trim(path, strDelim)
	}

	return strings.Trim(path[i+1:], strDelim)
}
