package formatter

import (
	"bytes"
	"testing"
)

func TestFriendly_printStatistics(t *testing.T) {
	tests := []struct {
		name     string
		stats    map[string]int
		expected string
	}{
		{
			name:     "no stats",
			stats:    map[string]int{},
			expected: "",
		},
		{
			name:     "nil stats",
			stats:    nil,
			expected: "",
		},
		{
			name:     "single stat",
			stats:    map[string]int{"rule1": 1},
			expected: "Warnings:\n  1  rule1  \n\n",
		},
		{
			name:     "multiple stats sorted by failures desc",
			stats:    map[string]int{"rule2": 2, "rule1": 1, "rule3": 3},
			expected: "Warnings:\n  3  rule3  \n  2  rule2  \n  1  rule1  \n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Friendly{}
			var buf bytes.Buffer
			f.printStatistics(&buf, "Warnings:", tt.stats)
			if got := buf.String(); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
