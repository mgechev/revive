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
