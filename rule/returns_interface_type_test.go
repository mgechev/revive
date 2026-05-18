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
	base := r.DefaultFilteredTypes()
	for _, v := range add {
		base[v] = struct{}{}
	}

	return base
}

func TestReturnsInterfaceTypeRule_Configure(t *testing.T) {
	tests := []struct {
		name               string
		arguments          lint.Arguments
		wantErr            error
		wantSearchingNames map[string]struct{}
	}{
		{
			name:               "no arguments",
			arguments:          lint.Arguments{},
			wantErr:            nil,
			wantSearchingNames: cfgInit([]string{}),
		},
		{
			name: "user defined searchingNames check list",
			arguments: lint.Arguments{
				map[string]any{
					"searchingNames": []any{
						"A",
					},
				},
			},
			wantErr:            nil,
			wantSearchingNames: cfgInit([]string{"A"}),
		},
		{
			name: "user defined searchingNames ok",
			arguments: lint.Arguments{
				map[string]any{
					"searchingNames": []any{
						"B",
					},
				},
			},
			wantErr:            nil,
			wantSearchingNames: cfgInit([]string{"B"}),
		},
		{
			name: "user defined searchingnames",
			arguments: lint.Arguments{
				map[string]any{
					"searchingnames": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr:            nil,
			wantSearchingNames: cfgInit([]string{"fixtures.DummyConfig"}),
		},
		{
			name: "user defined searching-names",
			arguments: lint.Arguments{
				map[string]any{
					"searching-names": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr:            nil,
			wantSearchingNames: cfgInit([]string{"fixtures.DummyConfig"}),
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
			name: "invalid user defined searching names",
			arguments: lint.Arguments{
				map[string]any{
					"searchingNames": []int{
						1,
					},
				},
			},
			wantErr: errors.New(`invalid format for entry 'searchingNames' of 'returns-interface-type' rule configuration: []string expected. got '[1]' ([]int)`),
		},
		{
			name: "invalid user defined searching names int values",
			arguments: lint.Arguments{
				map[string]any{
					"searchingNames": []any{
						1,
					},
				},
			},
			wantErr: errors.New(`invalid format for value in 'searchingNames' of 'returns-interface-type' rule configuration: string expected. got '1' (int)`),
		},
		{
			name: "user defined invalid searching-names key",
			arguments: lint.Arguments{
				map[string]any{
					"searching": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr: errors.New(`invalid argument 'searching' of 'returns-interface-type' rule configuration: searching-names expected. got 'searching'`),
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
			if !reflect.DeepEqual(rule.getFilteredTypes(), tt.wantSearchingNames) {
				t.Errorf("unexpected searchingNames: got = %v, want %v", rule.getFilteredTypes(), tt.wantSearchingNames)
			}
		})
	}
}
