package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringFormat(t *testing.T) {
	testRule(t, "string-format", &rule.StringFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]interface{}{
				"stringFormatMethod1", // The first argument is checked by default
				"/^[A-Z]/",
				"must start with a capital letter"},

			[]interface{}{
				"stringFormatMethod2[2].d",
				"/[^\\.]$/"}, // Must not end with a period
			[]interface{}{
				"s.Method3[2]",
				"!/^[Tt][Hh]/",
				"must not start with 'th'"},
			[]interface{}{
				"s.Method4", // same as before, but called from a struct
				"!/^[Ot][Tt]/",
				"must not start with 'ot'"}}})
}

func TestStringFormatArgumentParsing(t *testing.T) {
	r := &rule.StringFormatRule{}
	type argumentsTest struct {
		name          string
		config        lint.Arguments
		expectedError *string
	}
	stringPtr := func(s string) *string {
		return &s
	}
	tests := []argumentsTest{
		{
			name: "Not a Slice",
			config: lint.Arguments{
				"this is not a slice"},
			expectedError: stringPtr("invalid configuration for string-format: argument is not a slice [argument 0, option 0]")},
		{
			name: "Missing Regex",
			config: lint.Arguments{
				[]interface{}{
					"method[0]"}},
			expectedError: stringPtr("invalid configuration for string-format: less than two slices found in argument, scope and regex are required [argument 0, option 0]")},
		{
			name: "Bad Argument Type",
			config: lint.Arguments{
				[]interface{}{
					1}},
			expectedError: stringPtr("invalid configuration for string-format: less than two slices found in argument, scope and regex are required [argument 0, option 0]")},
		{
			name: "Empty Scope",
			config: lint.Arguments{
				[]interface{}{
					"",
					"//"}},
			expectedError: stringPtr("invalid configuration for string-format: empty scope provided [argument 0, option 0]")},
		{
			name: "Small or Empty Regex",
			config: lint.Arguments{
				[]interface{}{
					"method[1].a",
					"-"}},
			expectedError: stringPtr("invalid configuration for string-format: regex is too small (regexes should begin and end with '/') [argument 0, option 1]")},
		{
			name: "Bad Scope",
			config: lint.Arguments{
				[]interface{}{
					"1.a",
					"//"}},
			expectedError: stringPtr("failed to parse configuration for string-format: unable to parse rule scope [argument 0, option 0]")},
		{
			name: "Bad Regex",
			config: lint.Arguments{
				[]interface{}{
					"method[1].a",
					"/(/"}},
			expectedError: stringPtr("failed to parse configuration for string-format: unable to compile /(/ as regexp [argument 0, option 1]")},
		{
			name: "Sample Config",
			config: lint.Arguments{
				[]interface{}{
					"core.WriteError[1].Message", "/^([^A-Z]$)/", "must not start with a capital letter"},
				[]interface{}{
					"fmt.Errorf[0]", "/^|[^\\.!?]$/", "must not end in punctuation"},
				[]interface{}{
					"panic", "/^[^\\n]*$/", "must not contain line breaks"}}},
		{
			name: "Underscores in Scope",
			config: lint.Arguments{
				[]interface{}{
					"some_pkg._some_function_name[5].some_member",
					"//"}}}}

	for _, a := range tests {
		t.Run(a.name, func(t *testing.T) {
			err := r.ParseArgumentsTest(a.config)
			if err != nil {
				if a.expectedError == nil || *err != *a.expectedError {
					t.Errorf("unexpected panic message: %s", *err)
				}
			} else if a.expectedError != nil {
				t.Error("error expected but not received")
			}
		})
	}
}
