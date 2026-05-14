package rule

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestReturnsInterfaceTypeRule_Configure(t *testing.T) {
	tests := []struct {
		name                        string
		arguments                   lint.Arguments
		wantErr                     error
		wantStopOnFirst             bool
		wantUserDefinedIgnoredNames map[string]struct{}
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
			wantErr:   nil,
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
			name: "invalid stopOnFirst value",
			arguments: lint.Arguments{
				map[string]any{
					"stopOnFirst": "abc",
				},
			},
			wantErr: errors.New(`invalid argument 'stopOnFirst' for 'returns-interface-type' rule, expecting bool value. Got 'abc' (string)`),
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
			wantErr: errors.New(`invalid format 'ignoredNames' for 'returns-interface-type' rule []string expected. Got '[1]' ([]int)`),
		},
		{
			name: "invalid user defined ignored names int values",
			arguments: lint.Arguments{
				map[string]any{
					"stopOnFirst": false,
					"ignoredNames": []any{
						1,
					},
				},
			},
			wantErr: errors.New(`invalid value in 'ignoredNames' for 'returns-interface-type' rule string expected Got '1' (int)`),
		},

		{
			name: "user defined ignored names",
			arguments: lint.Arguments{
				map[string]any{
					"stopOnFirst": false,
					"ignoredNames": []any{
						"fixtures.DummyConfig",
					},
				},
			},
			wantErr: nil,
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
			if rule.stopOnFirst != tt.wantStopOnFirst {
				t.Errorf("unexpected stopOnFirst: got = %v, want %v", rule.stopOnFirst, tt.wantStopOnFirst)
			}
		})
	}
}

func TestReturnsInterfaceTypeRule_Configure_LoadTypes(t *testing.T) {
	t.Run("loads types", func(t *testing.T) {
		var rule ReturnsInterfaceTypeRule

		err := rule.Configure(lint.Arguments{
			map[string]any{
				"stopOnFirst": false,
				"ignoredNames": []any{
					"fixtures.DummyResults",
				},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: got = %v, want = nil", err)
		}

		types := []string{"fixtures.DummyResults"}
		for _, typeValue := range types {
			all := rule.getIgnoredTypes()
			_, ok := all[typeValue]
			if !ok {
				t.Errorf("not loaded expected type %q", typeValue)
			}
		}
	})
}
