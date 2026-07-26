package rule

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"github.com/mgechev/revive/lint"
)

// PrimitiveInNameRule lints the name of a variable.
type PrimitiveInNameRule struct {
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *PrimitiveInNameRule) Configure(_ lint.Arguments) error {
	return nil
}

// Apply applies the rule to given file.
func (*PrimitiveInNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	check := func(id *ast.Ident) {
		if id.Name == "_" {
			return
		}

		word, ok := isPrimitiveResolvePrimitiveString(file.Pkg.TypeOf(id))
		if !ok || !hasWord(id.Name, word) {
			return
		}

		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryNaming,
			Confidence: 0.8,
			Node:       id,
			Failure:    "avoid primitive type in name",
		})
	}

	ast.Inspect(file.AST, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.AssignStmt:
			if v.Tok != token.DEFINE {
				return true
			}
			for _, expr := range v.Lhs {
				if id, ok := expr.(*ast.Ident); ok {
					check(id)
				}
			}
		case *ast.GenDecl:
			if v.Tok != token.VAR && v.Tok != token.CONST {
				return true
			}
			for _, spec := range v.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, id := range vs.Names {
					check(id)
				}
			}
		}
		return true
	})

	return failures
}

// isPrimitiveResolvePrimitiveString will resolve the type to a string.
// Note currently this will not work for []string, only concrete types.
func isPrimitiveResolvePrimitiveString(typ types.Type) (string, bool) {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return "", false
	}

	switch basic.Kind() {
	case types.Int, types.Int8, types.Int16, types.Int64,
		types.Uint, types.Uint16, types.Uint32, types.Uint64, types.Uintptr,
		types.UntypedInt:
		return "Int", true
	case types.Int32, types.UntypedRune:
		return "Rune", true
	case types.Uint8:
		return "Byte", true
	case types.String, types.UntypedString:
		return "String", true
	case types.Bool, types.UntypedBool:
		return "Bool", true
	case types.Float32:
		return "Float32", true
	case types.Float64, types.UntypedFloat:
		return "Float64", true
	}

	return "", false
}

func hasWord(name, word string) bool {
	for _, segment := range splitWords(name) {
		if strings.EqualFold(segment, word) {
			return true
		}
	}
	return false
}

func splitWords(name string) []string {
	var words []string
	var current []rune

	runes := []rune(name)
	for i, r := range runes {
		if r == '_' {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
			continue
		}
		if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(runes[i-1]) {
			words = append(words, string(current))
			current = nil
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}

// Name returns the rule name.
func (*PrimitiveInNameRule) Name() string {
	return "primitive-in-name"
}
