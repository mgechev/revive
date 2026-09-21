package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDotImportsRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
		},
		{
			name: "no allowedPackages key",
			arguments: lint.Arguments{
				map[string]any{
					"invalid": "argument",
				},
			},
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowedPackages": []any{
						"github.com/onsi/ginkgo/v2",
					},
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowedpackages": []any{
						"github.com/onsi/ginkgo/v2",
					},
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowed-packages": []any{
						"github.com/onsi/ginkgo/v2",
					},
				},
			},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument to the dot-imports rule. Expecting a k,v map, got string"),
		},
		{
			name: "invalid allowedPackages type",
			arguments: lint.Arguments{
				map[string]any{
					"allowedPackages": "invalid",
				},
			},
			wantErr: errors.New("invalid argument to the dot-imports rule, []string expected. Got 'invalid' (string)"),
		},
		{
			name: "invalid allowedPackages value type",
			arguments: lint.Arguments{
				map[string]any{
					"allowedPackages": []any{123},
				},
			},
			wantErr: errors.New("invalid argument to the dot-imports rule, string expected. Got '123' (int)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.DotImportsRule

			err := r.Configure(tt.arguments)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Configure() unexpected non-nil error %q", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr.Error() {
				t.Errorf("Configure() unexpected error: got %q, want %q", err, tt.wantErr)
			}
		})
	}
}
