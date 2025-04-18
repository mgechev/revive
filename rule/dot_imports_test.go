package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestDotImportsRule_Configure(t *testing.T) {
	tests := []struct {
		name              string
		arguments         lint.Arguments
		wantErr           error
		wantAllowPackages allowPackages
	}{
		{
			name:              "no arguments",
			arguments:         lint.Arguments{},
			wantErr:           nil,
			wantAllowPackages: allowPackages{},
		},
		{
			name: "no allowedPackages key",
			arguments: lint.Arguments{
				map[string]any{
					"invalid": "argument",
				},
			},
			wantErr:           nil,
			wantAllowPackages: allowPackages{},
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
			wantErr: nil,
			wantAllowPackages: allowPackages{
				`"github.com/onsi/ginkgo/v2"`: struct{}{},
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
			wantErr: nil,
			wantAllowPackages: allowPackages{
				`"github.com/onsi/ginkgo/v2"`: struct{}{},
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
			wantErr: nil,
			wantAllowPackages: allowPackages{
				`"github.com/onsi/ginkgo/v2"`: struct{}{},
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
			var rule DotImportsRule

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
			if !reflect.DeepEqual(rule.allowedPackages, tt.wantAllowPackages) {
				t.Errorf("unexpected allowedPackages: got = %v, want %v", rule.allowedPackages, tt.wantAllowPackages)
			}
		})
	}
}
