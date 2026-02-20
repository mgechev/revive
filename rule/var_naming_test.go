package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestVarNamingRule_Configure(t *testing.T) {
	tests := []struct {
		name                         string
		arguments                    lint.Arguments
		wantErr                      error
		wantAllowList                []string
		wantBlockList                []string
		wantSkipInitialismNameChecks bool
		wantAllowUpperCaseConst      bool
	}{
		{
			name:                         "no arguments",
			arguments:                    lint.Arguments{},
			wantErr:                      nil,
			wantAllowList:                nil,
			wantBlockList:                nil,
			wantSkipInitialismNameChecks: false,
			wantAllowUpperCaseConst:      false,
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skipInitialismNameChecks": true,
					"upperCaseConst":           true,
				}},
			},
			wantErr:                      nil,
			wantAllowList:                []string{"ID"},
			wantBlockList:                []string{"VM"},
			wantSkipInitialismNameChecks: true,
			wantAllowUpperCaseConst:      true,
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skipinitialismnamechecks": true,
					"uppercaseconst":           true,
				}},
			},
			wantErr:                      nil,
			wantAllowList:                []string{"ID"},
			wantBlockList:                []string{"VM"},
			wantSkipInitialismNameChecks: true,
			wantAllowUpperCaseConst:      true,
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skip-initialism-name-checks": true,
					"upper-case-const":            true,
				}},
			},
			wantErr:                      nil,
			wantAllowList:                []string{"ID"},
			wantBlockList:                []string{"VM"},
			wantSkipInitialismNameChecks: true,
			wantAllowUpperCaseConst:      true,
		},
		{
			name:      "invalid allowlist type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid allowlist value type",
			arguments: lint.Arguments{[]any{123}},
			wantErr:   errors.New("invalid 123 values of the var-naming rule. Expecting slice of strings but got element of type []interface {}"),
		},
		{
			name:      "invalid blocklist type",
			arguments: lint.Arguments{[]any{"ID"}, 123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a blocklist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid third argument type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, 123},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, got int"),
		},
		{
			name:      "invalid third argument slice size",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, but 0"),
		},
		{
			name:      "invalid third argument first element type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{123}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, with map, but int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule VarNamingRule

			err := rule.Configure(tt.arguments)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("unexpected error: got = nil, want = %v", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("unexpected error: got = %v, want = %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: got = %v, want = nil", err)
			}
			if rule.allowUpperCaseConst != tt.wantAllowUpperCaseConst {
				t.Errorf("unexpected allowUpperCaseConst: got = %v, want %v", rule.allowUpperCaseConst, tt.wantAllowUpperCaseConst)
			}
			if rule.skipInitialismNameChecks != tt.wantSkipInitialismNameChecks {
				t.Errorf("unexpected skipInitialismNameChecks: got = %v, want %v", rule.skipInitialismNameChecks, tt.wantSkipInitialismNameChecks)
			}
		})
	}
}

func TestHasUpperCaseLetter(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", true},
		{"fmt", false},
		{"_SOME_PRIVATE_CONST_2", true},
		{"", false},
		// Unicode uppercase (non-ASCII)
		{"Ä", false}, // Latin capital letter A with diaeresis
		{"Ω", false}, // Greek capital letter Omega
		{"Д", false}, // Cyrillic capital letter De

		// Unicode lowercase/symbols
		{"ß", false}, // German sharp s
		{"π", false}, // Greek small letter pi
		{"💡", false}, // Emoji
		{"你", false}, // Chinese character
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := hasUpperCaseLetter(tt.varName); got != tt.expected {
				t.Errorf("hasCaps(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func TestIsUpperCaseConst(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"SOME_CONST_2", true},
		{"__FOO", false},
		{"__", false},
		{"X509B", true},
		{"FOO", true},
		{"1FOO", false},
		{"_FOO123_BAR456", true},
		{"A1_B2_C3", true},
		{"A1_b2", false},
		{"FOO_", false},
		{"foo", false},
		{"_", false},
		{"", false},
		{"FOOBAR", true},
		{"FO", true},
		{"F_O", true},
		{"FOO123", true},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := isUpperCaseConst(tt.varName); got != tt.expected {
				t.Errorf("isUpperCaseConst(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func TestIsUpperUnderscore(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"_", false},
		{"", false},
		{"empty string", false},
		{"_404_404", true},
		{"FOO_BAR", true},
		{"FOOBAR", false},
		{"FO", false},
		{"F_O", false},
		{"_FOOBAR", true},
		{"FOOBAR_", true},
		{"FOO123", false},
		{"FOO_123", true},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := isUpperUnderscore(tt.varName); got != tt.expected {
				t.Errorf("isUpperUnderScore(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{'0', true},
		{'1', true},
		{'2', true},
		{'3', true},
		{'4', true},
		{'5', true},
		{'6', true},
		{'7', true},
		{'8', true},
		{'9', true},
		{'a', false},
		{'Z', false},
		{' ', false},
		{'!', false},
		{'🙂', false}, // Emoji to test unicode
		{'٠', false}, // Arabic-Indic 0
		{'١', false}, // Arabic-Indic 1
		{'२', false}, // Devanagari 2
		{'৩', false}, // Bengali 3
		{'४', false}, // Devanagari 4
		{'௫', false}, // Tamil 5
		{'๖', false}, // Thai 6
		{'৭', false}, // Bengali 7
		{'८', false}, // Devanagari 8
		{'९', false}, // Devanagari 9
	}

	for _, tt := range tests {
		result := isDigit(tt.input)
		if result != tt.expected {
			t.Errorf("isDigit(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsUpper(t *testing.T) {
	t.Run("non letter", func(t *testing.T) {
		tests := []rune{
			'0',
			'5',
			' ',
			'_',
			'!',
			'🙂', // Emoji to test unicode
		}
		for _, r := range tests {
			result := isUpper(r)
			if result {
				t.Errorf("isUpper(%q) = %v; want false", r, result)
			}
		}
	})

	t.Run("non ASCII letter", func(t *testing.T) {
		tests := []rune{
			'Ą',
			'Ć',
			'你',
			'日',
			'本',
			'語',
			'韓',
			'中',
			'文',
			'あ',
			'ア',
			'한',
		}
		for _, r := range tests {
			result := isUpper(r)
			if result {
				t.Errorf("isUpper(%q) = %v; want false", r, result)
			}
		}
	})

	t.Run("lowercase ASCII letter", func(t *testing.T) {
		tests := []rune{
			'a',
			'b',
		}
		for _, r := range tests {
			result := isUpper(r)
			if result {
				t.Errorf("isUpper(%q) = %v; want false", r, result)
			}
		}
	})

	t.Run("uppercase ASCII letter", func(t *testing.T) {
		tests := []rune{
			'A',
			'B',
			'C',
			'Z',
		}
		for _, r := range tests {
			result := isUpper(r)
			if !result {
				t.Errorf("isUpper(%q) = %v; want true", r, result)
			}
		}
	})
}
