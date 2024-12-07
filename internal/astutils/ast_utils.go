package astutils

import "go/ast"

// FuncSignatureIs returns true if the given func decl satisfies a signature characterized
// by the given name, parameters types and return types; false otherwise.
//
// Example: to check if a function declaration has the signature Foo(int, string) (bool,error)
// call to FuncSignatureIs(funcDecl,"Foo",[]string{"int","string"},[]string{"bool","error"})
func FuncSignatureIs(funcDecl *ast.FuncDecl, wantName string, wantParametersTypes, wantResultsTypes []string) bool {
	if wantName != funcDecl.Name.String() {
		return false // func name doesn't match expected one
	}

	funcParametersTypes := getTypeNames(funcDecl.Type.Params)
	if len(wantParametersTypes) != len(funcParametersTypes) {
		return false // func has not the expected number of parameters
	}

	funcResultsTypes := getTypeNames(funcDecl.Type.Results)
	if len(wantResultsTypes) != len(funcResultsTypes) {
		return false // func has not the expected number of return values
	}

	for i, wantType := range wantParametersTypes {
		if wantType != funcParametersTypes[i] {
			return false // type of a func's parameter does not match the type of the corresponding expected parameter
		}
	}

	for i, wantType := range wantResultsTypes {
		if wantType != funcResultsTypes[i] {
			return false // type of a func's return value does not match the type of the corresponding expected return value
		}
	}

	return true
}

func getTypeNames(fields *ast.FieldList) []string {
	result := []string{}

	if fields == nil {
		return result
	}

	for _, field := range fields.List {
		typeName := field.Type.(*ast.Ident).Name
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
