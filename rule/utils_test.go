package rule

import (
	"fmt"
	"testing"
)

func TestIsCallToExitFunction(t *testing.T) {
	tests := []struct {
		pkgName      string
		functionName string
		expected     bool
	}{
		{"os", "Exit", true},
		{"syscall", "Exit", true},
		{"log", "Fatal", true},
		{"log", "Fatalf", true},
		{"log", "Fatalln", true},
		{"log", "Panic", true},
		{"log", "Panicf", true},
		{"log", "Print", false},
		{"fmt", "Errorf", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.pkgName, tt.functionName), func(t *testing.T) {
			if got := isCallToExitFunction(tt.pkgName, tt.functionName); got != tt.expected {
				t.Errorf("isCallToExitFunction(%s, %s) = %v; want %v", tt.pkgName, tt.functionName, got, tt.expected)
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
