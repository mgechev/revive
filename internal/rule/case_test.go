package rule_test

import (
	"testing"

	"github.com/mgechev/revive/internal/rule"
)

func TestHasUpperCaseLetter(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", true},
		{"fmt", false},
		{"_SOME_PRIVATE_CONST_2", true},
		{"", false},
		{"a1_ b!", false},
		// Unicode uppercase (non-ASCII)
		{"Ä", false}, // Latin capital letter A with diaeresis
		{"Ą", false}, // Latin capital letter A with ogonek
		{"Ω", false}, // Greek capital letter Omega
		{"Д", false}, // Cyrillic capital letter De

		// Unicode lowercase/symbols
		{"ß", false}, // German sharp s
		{"π", false}, // Greek small letter pi
		{"💡", false}, // Emoji
		{"你", false}, // Chinese character
		{"日本語", false},
		{"あア", false},
		{"한", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := rule.HasUpperCaseLetter(tt.varName); got != tt.expected {
				t.Errorf("HasUpperCaseLetter(%s) = %v; want %v", tt.varName, got, tt.expected)
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
		{"AB", true},
		{"F_O", true},
		{"FOO123", true},
		{"FOO0123456789", true},
		{"FOO BAR", false},
		{"FOO!", false},
		{"FOO🙂", false},
		// Only ASCII digits are accepted
		{"FOO٠", false}, // Arabic-Indic 0
		{"FOO١", false}, // Arabic-Indic 1
		{"FOO२", false}, // Devanagari 2
		{"FOO৩", false}, // Bengali 3
		{"FOO४", false}, // Devanagari 4
		{"FOO௫", false}, // Tamil 5
		{"FOO๖", false}, // Thai 6
		{"FOO৭", false}, // Bengali 7
		{"FOO८", false}, // Devanagari 8
		{"FOO९", false}, // Devanagari 9
		// Only ASCII uppercase letters are accepted
		{"FOOĄ", false},
		{"FOOĆ", false},
		{"FOO你", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := rule.IsUpperCaseConst(tt.varName); got != tt.expected {
				t.Errorf("IsUpperCaseConst(%s) = %v; want %v", tt.varName, got, tt.expected)
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
		{"AB", false},
		{"F_O", false},
		{"_FOOBAR", true},
		{"FOOBAR_", true},
		{"FOO123", false},
		{"FOO_123", true},
		{"FOO_BAR!", false},
		{"FOO_BAR٠", false}, // Arabic-Indic 0
		{"FOO_BARĄ", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := rule.IsUpperUnderscore(tt.varName); got != tt.expected {
				t.Errorf("IsUpperUnderscore(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}
