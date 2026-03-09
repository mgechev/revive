package rule

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"
)

func TestGetImportPackageNameUsesImportedPackageNameWhenTypeInfoAvailable(t *testing.T) {
	alias := ast.NewIdent("utls")
	imp := &ast.ImportSpec{
		Name: alias,
		Path: &ast.BasicLit{Value: `"github.com/enetx/utls"`},
	}

	typesInfo := &types.Info{Defs: map[*ast.Ident]types.Object{}}
	typesInfo.Defs[alias] = types.NewPkgName(token.NoPos, nil, alias.Name, types.NewPackage("github.com/enetx/utls", "tls"))

	got := getImportPackageName(imp, typesInfo)
	if got != "tls" {
		t.Fatalf("expected imported package name %q, got %q", "tls", got)
	}
}
