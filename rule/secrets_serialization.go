package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mgechev/revive/lint"
)

var defaultSecretFieldIndicators = []string{
	"BearerToken", "Secret", "Token", "Password", "Key", "APIKey", "Auth", "Credential", "ClientSecret", "AccessToken", "AuthToken",
}

type SecretsSerializationRule struct {
	secretFieldIndicators map[string]struct{}
}

func (r *SecretsSerializationRule) Name() string {
	return "secrets-serialization"
}

func (r *SecretsSerializationRule) Configure(arguments lint.Arguments) error {
	var err error
	indicators := []string{}
	if len(arguments) < 1 {
		indicators = defaultSecretFieldIndicators
	} else {
		indicators, err = r.getSecretFieldIndicatorList(arguments[0], "secretFieldIndicators")
	}
	for _, indicator := range indicators {
		r.secretFieldIndicators[strings.ToLower(indicator)] = struct{}{}
	}
	return err
}

func (r *SecretsSerializationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	walker := func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				for _, fieldName := range field.Names {
					if !ast.IsExported(fieldName.Name) {
						continue
					}
					if field.Tag == nil || strings.Contains(field.Tag.Value, `json:"-"`) {
						continue
					}
					if !r.isLikelySecret(fieldName.Name) {
						continue
					}
					failures = append(failures, lint.Failure{
						Confidence: 0.8,
						Node:       field,
						Category:   lint.FailureCategoryBadPractice,
						Failure:    "Struct field '" + fieldName.Name + "' may contain secrets but is not excluded from JSON serialization (missing `json:\"-\"`)",
					})
				}
			}
		}
		return true
	}
	ast.Inspect(file.AST, walker)
	return failures
}

func (r *SecretsSerializationRule) isLikelySecret(name string) bool {
	_, ok := r.secretFieldIndicators[strings.ToLower(name)]
	return ok
}

func (r *SecretsSerializationRule) getSecretFieldIndicatorList(arg any, argName string) ([]string, error) {
	args, ok := arg.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid argument to the secrets-serialization rule: expecting %s of type slice of strings, got %T", argName, arg)
	}
	var list []string
	for _, v := range args {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid argument to the secrets-serialization rule: expecting %s of type slice of strings, got slice of type %T", argName, v)
		}
		list = append(list, val)
	}
	return list, nil
}
