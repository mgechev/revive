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

func TestAnyCapsFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", true},
		{"fmt", false},
		{"_SOME_PRIVATE_CONST_2", true},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := anyCaps(tt.varName); got != tt.expected {
				t.Errorf("anyCaps(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func BenchmarkAnyCapsRE(b *testing.B) {

	var anyCapsRE = regexp.MustCompile(`[A-Z]`)
	for i := 0; i <= b.N; i++ {
		input := "HeLlo_WoRlD"
		_ = anyCapsRE.MatchString(input)
	}
}

func BenchmarkAnyCaps(b *testing.B) {
	for i := 0; i <= b.N; i++ {		
		input := "HeLlo_WoRlD"
		_ = anyCaps(input)
	}
}

func TestAllCapsFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
		{"Exit", false},
		{"fmt", false},
		{"_SOME_PRIVATE_CONST_2", true},
		{"HELLO_WORLD123", true},
		{"Hello_World", false},
		{"", false},
		{"INVALID-CHAR", false},
	}

	for _, tt := range tests {
		t.Run(tt.varName, func(t *testing.T) {
			if got := allCaps(tt.varName); got != tt.expected {
				t.Errorf("allCaps(%s) = %v; want %v", tt.varName, got, tt.expected)
			}
		})
	}
}

func BenchmarkAllCapsRE(b *testing.B) {

	var allCapsRE = regexp.MustCompile(`^_?[A-Z][A-Z\d]*(_[A-Z\d]+)*$`)
	for i := 0; i <= b.N; i++ {
		input := "_SOME_PRIVATE_CONST_2"
		_ = allCapsRE.MatchString(input)
	}
}

func BenchmarkAllCaps(b *testing.B) {
	for i := 0; i <= b.N; i++ {	
		input := "_SOME_PRIVATE_CONST_2"
		_ = allCaps(input)
	}
}


func TestIsUpperConstFunction(t *testing.T) {
	tests := []struct {
		varName  string
		expected bool
	}{
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
