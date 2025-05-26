// Package astutils provides utility functions for working with AST nodes
package astutils

import (
	"go/ast"
	"go/token"
	"slices"
)

// FuncSignatureIs returns true if the given func decl satisfies a signature characterized
// by the given name, parameters types and return types; false otherwise.
//
// Example: to check if a function declaration has the signature Foo(int, string) (bool,error)
// call to FuncSignatureIs(funcDecl,"Foo",[]string{"int","string"},[]string{"bool","error"})
func FuncSignatureIs(funcDecl *ast.FuncDecl, wantName string, wantParametersTypes, wantResultsTypes []string) bool {
	if wantName != funcDecl.Name.String() {
		return false // func name doesn't match expected one
	}

	funcResultsTypes := GetTypeNames(funcDecl.Type.Results)
	if !slices.Equal(wantResultsTypes, funcResultsTypes) {
		return false // func has not the expected return values
	}

	// Name and return values are those we expected,
	// the final result depends on parameters being what we want.
	return funcParametersSignatureIs(funcDecl, wantParametersTypes)
}

// funcParametersSignatureIs returns true if the function has parameters of the given type and order,
// false otherwise
func funcParametersSignatureIs(funcDecl *ast.FuncDecl, wantParametersTypes []string) bool {
	funcParametersTypes := GetTypeNames(funcDecl.Type.Params)

	return slices.Equal(wantParametersTypes, funcParametersTypes)
}

// GetTypeNames yields an slice with the string representation of the types of given fields.
// It yields nil if the field list is nil.
func GetTypeNames(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}

	result := []string{}

	for _, field := range fields.List {
		typeName := getFieldTypeName(field.Type)
		if field.Names == nil { // unnamed field
			result = append(result, typeName)
			continue
		}

		for range field.Names { // add one type name for each field name
			result = append(result, typeName)
		}
	}

	return result
}

func getFieldTypeName(typ ast.Expr) string {
	switch f := typ.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return getFieldTypeName(f.X) + "." + getFieldTypeName(f.Sel)
	case *ast.StarExpr:
		return "*" + getFieldTypeName(f.X)
	case *ast.IndexExpr:
		return getFieldTypeName(f.X) + "[" + getFieldTypeName(f.Index) + "]"
	case *ast.ArrayType:
		return "[]" + getFieldTypeName(f.Elt)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "UNHANDLED_TYPE"
	}
}

// IsStringLiteral returns true if the given expression is a string literal, false otherwise
func IsStringLiteral(e ast.Expr) bool {
	sl, ok := e.(*ast.BasicLit)

	return ok && sl.Kind == token.STRING
}
