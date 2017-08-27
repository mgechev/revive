package visitors

import (
	"errors"
	"reflect"

	"github.com/mgechev/golinter/file"
	"github.com/mgechev/golinter/rules"
)

// Setup sets the proper pointers of given visitor.
func Setup(v interface{}, conf rules.RuleConfig, file *file.File) error {
	val := reflect.ValueOf(v).Elem()
	field := val.FieldByName("RuleVisitor")
	if !field.IsValid() {
		return errors.New("invalid rule visitor")
	}
	field.Set(reflect.ValueOf(RuleVisitor{RuleName: conf.Name, RuleArguments: conf.Arguments, File: file}))

	field = val.FieldByName("Impl")
	if !field.IsValid() {
		return errors.New("invalid rule visitor")
	}
	field.Set(reflect.ValueOf(v))

	return nil
}
