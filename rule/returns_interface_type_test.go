package rule

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
)

// cfgInit helper for create test config with default + case update.
func cfgInit(add []string) map[string]struct{} {
	var r ReturnsInterfaceTypeRule
	base := r.DefaultIgnoredTypes()
	for _, v := range add {
		base[v] = struct{}{}
	}

	return base
}

func TestReturnsInterfaceTypeRule_Configure(t *testing.T) {
	tests := []struct {
		name             string
		arguments        lint.Arguments
		wantErr          error
		wantIgnoredNames map[string]struct{}
	}{
		{
			name:             "no arguments",
			arguments:        lint.Arguments{},
			wantErr:          nil,
			wantIgnoredNames: cfgInit([]string{}),
		},
		{
			name: "user defined ignoredNames check list",
			arguments: lint.Arguments{
				map[string]any{
					"ignoredNames": []any{
						"A",
					},
				},
			},
			wantErr:          nil,
			wantIgnoredNames: cfgInit([]string{"A"}),
		},
		{
			name: "user defined ignoredNames ok",
			arguments: lint.Arguments{
				map[string]any{
					"ignoredNames": []any{
						"B",
					},
				},
			},
			wantErr:          nil,
			wantIgnoredNames: cfgInit([]string{"B"}),
		},
		{
			name: "user defined ignorednames",
			arguments: lint.Arguments{
				map[string]any{
					"ignorednames": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr:          nil,
			wantIgnoredNames: cfgInit([]string{"fixtures.DummyConfig"}),
		},
		{
			name: "user defined ignored-names",
			arguments: lint.Arguments{
				map[string]any{
					"ignored-names": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr:          nil,
			wantIgnoredNames: cfgInit([]string{"fixtures.DummyConfig"}),
		},

		{
			name: "invalid arguments format",
			arguments: lint.Arguments{
				[]any{
					"abc",
				},
			},
			wantErr: errors.New(`invalid argument '[abc]' for 'returns-interface-type' rule. Expecting a k,v map, got []interface {}`),
		},
		{
			name: "invalid user defined ignored names",
			arguments: lint.Arguments{
				map[string]any{
					"ignoredNames": []int{
						1,
					},
				},
			},
			wantErr: errors.New(`invalid format for entry 'ignoredNames' of 'returns-interface-type' rule configuration: []string expected. got '[1]' ([]int)`),
		},
		{
			name: "invalid user defined ignored names int values",
			arguments: lint.Arguments{
				map[string]any{
					"ignoredNames": []any{
						1,
					},
				},
			},
			wantErr: errors.New(`invalid format for value in 'ignoredNames' of 'returns-interface-type' rule configuration: string expected. got '1' (int)`),
		},
		{
			name: "user defined invalid ignored-names key",
			arguments: lint.Arguments{
				map[string]any{
					"ignored": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr: errors.New(`invalid argument 'ignored' of 'returns-interface-type' rule configuration: ignored-names expected. got 'ignored'`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule ReturnsInterfaceTypeRule

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
			if !reflect.DeepEqual(rule.getIgnoredTypes(), tt.wantIgnoredNames) {
				t.Errorf("unexpected ignoredNames: got = %v, want %v", rule.getIgnoredTypes(), tt.wantIgnoredNames)
			}
		})
	}
}
