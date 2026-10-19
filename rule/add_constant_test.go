package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestAddConstantRule_Configure(t *testing.T) {
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
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowFloats": "1.0,2.0",
					"allowInts":   "1,2",
					"allowStrs":   "a,b",
					"maxLitCount": "3",
					"ignoreFuncs": "fmt.Println,fmt.Printf",
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allowfloats": "1.0,2.0",
					"allowints":   "1,2",
					"allowstrs":   "a,b",
					"maxlitcount": "3",
					"ignorefuncs": "fmt.Println,fmt.Printf",
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"allow-floats":  "1.0,2.0",
					"allow-ints":    "1,2",
					"allow-strs":    "a,b",
					"max-lit-count": "3",
					"ignore-funcs":  "fmt.Println,fmt.Printf",
				},
			},
		},
		{
			name: "unrecognized key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "someValue",
				},
			},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument to the add-constant rule, expecting a k,v map. Got string"),
		},
		{
			name: "invalid allowFloats value",
			arguments: lint.Arguments{
				map[string]any{
					"allowFloats": 123,
				},
			},
			wantErr: errors.New("invalid argument to the add-constant rule, string expected. Got '123' (int)"),
		},
		{
			name: "invalid maxLitCount value: not a string",
			arguments: lint.Arguments{
				map[string]any{
					"maxLitCount": 123,
				},
			},
			wantErr: errors.New("invalid argument to the add-constant rule, expecting string representation of an integer. Got '123' (int)"),
		},
		{
			name: "invalid maxLitCount value: not an int",
			arguments: lint.Arguments{
				map[string]any{
					"maxLitCount": "abc",
				},
			},
			wantErr: errors.New("invalid argument to the add-constant rule, expecting string representation of an integer. Got 'abc'"),
		},
		{
			name: "invalid ignoreFuncs value: not a string",
			arguments: lint.Arguments{
				map[string]any{
					"ignoreFuncs": 123,
				},
			},
			wantErr: errors.New("invalid argument to the ignoreFuncs parameter of add-constant rule, string expected. Got '123' (int)"),
		},
		{
			name: "invalid ignoreFuncs value: empty string",
			arguments: lint.Arguments{
				map[string]any{
					"ignoreFuncs": " ",
				},
			},
			wantErr: errors.New("invalid argument to the ignoreFuncs parameter of add-constant rule, expected regular expression must not be empty"),
		},
		{
			name: "invalid ignoreFuncs value: wrong regexp",
			arguments: lint.Arguments{
				map[string]any{
					"ignoreFuncs": "(",
				},
			},
			wantErr: errors.New(`invalid argument to the ignoreFuncs parameter of add-constant rule: regexp "(" does not compile: error parsing regexp: missing closing ): ` + "`(`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.AddConstantRule

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
