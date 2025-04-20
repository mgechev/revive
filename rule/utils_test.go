package rule

import (
	"fmt"
	"regexp"
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

func TestHasUpperCaseFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", true},
		{"fmt", false},
		{"_SOME_PRIVATE_CONST_2", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := hasUpperCase(tt.varName); got != tt.expected {
				t.Errorf("hasCaps(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func BenchmarkHasUpperCaseRE(b *testing.B) {
	var anyUpperCaseRE = regexp.MustCompile(`[A-Z]`)
	for i := 0; i <= b.N; i++ {
		input := "HeLlo_WoRlD"
		_ = anyUpperCaseRE.MatchString(input)
	}
}

func BenchmarkHasUpperCase(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		input := "HeLlo_WoRlD"
		_ = hasUpperCase(input)
	}
}

func TestHasLowerFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", true},
		{"fmt", true},
		{"_SOME_PRIVATE_CONST_2", false},
		{"HELLO_WORLD123", false},
		{"Hello_World", true},
		{"", false},
		{"INVALID-CHAR", false},
		{"abc123", true},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := hasLower(tt.varName); got != tt.expected {
				t.Errorf("hasLower(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func BenchmarkAllCapsRE(b *testing.B) {
	var allUpperCaseRE = regexp.MustCompile(`^_?[A-Z][A-Z\d]*(_[A-Z\d]+)*$`)
	for i := 0; i <= b.N; i++ {
		input := "_SOME_PRIVATE_CONST_2"
		_ = allUpperCaseRE.MatchString(input)
	}
}

func BenchmarkAllCaps(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		input := "_SOME_PRIVATE_CONST_2"
		_ = hasUpperCase(input)
	}
}

func TestIsUpperConstFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{

		{"X509B", true},
		{"FOO", true},
		{"_FOO123_BAR456", true},
		{"A1_B2_C3", true},
		{"A1_b2", false},
		{"__FOO", false},
		{"FOO_", false},
		{"foo", false},
		{"_", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := isUpperCaseConst(tt.varName); got != tt.expected {
				t.Errorf("isUpperCaseConst(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func BenchmarkUpperCaseConstRE(b *testing.B) {
	var upperCaseConstRE = regexp.MustCompile(`^_?[A-Z][A-Z\d]*(_[A-Z\d]+)*$`)
	for i := 0; i <= b.N; i++ {
		input := "A1_B2_C3"
		_ = upperCaseConstRE.MatchString(input)
	}
}

func BenchmarkIsUpperCaseConst(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		input := "A1_B2_C3"
		_ = isUpperCaseConst(input)
	}
}
