package structtag_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mgechev/revive/internal/structtag"
)

func TestParse(t *testing.T) {
	test := []struct {
		name    string
		tag     string
		want    []*structtag.Tag
		wantErr error
	}{
		{
			name: "empty tag",
			tag:  "",
			want: nil,
		},
		{
			name: "tag with only whitespace",
			tag:  "   ",
			want: nil,
		},
		{
			name: "tag with leading tabs and newlines",
			tag:  "\tjson:\"foo\"\n\txml:\"bar\"",
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "xml",
					Name: "bar",
				},
			},
		},
		{
			name: "tag with one key (valid)",
			tag:  `json:""`,
			want: []*structtag.Tag{
				{
					Key: "json",
				},
			},
		},
		{
			name: "tag with one key and dash name",
			tag:  `json:"-"`,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "-",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"foo"`,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with key, name and option",
			tag:  `json:"foo,omitempty"`,
			want: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
			},
		},
		{
			name: "tag with multiple keys",
			tag:  `json:"" hcl:""`,
			want: []*structtag.Tag{
				{
					Key: "json",
				},
				{
					Key: "hcl",
				},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" hcl:"foo"`,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with multiple keys and different names",
			tag:  `json:"foo" hcl:"bar"`,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "bar",
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo,omitempty" structs:"bar,omitnested"`,
			want: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo" structs:"bar,omitnested" hcl:"-"`,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
				{
					Key:  "hcl",
					Name: "-",
				},
			},
		},
		{
			name: "tag with quoted name",
			tag:  `json:"foo,bar:\"baz\""`,
			want: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{`bar:"baz"`},
				},
			},
		},
		{
			name: "tag with trailing space",
			tag:  `json:"foo" `,
			want: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
		{
			name:    "tag with one key (invalid)",
			tag:     "json",
			wantErr: errors.New("invalid syntax for struct tag pair"),
		},
		{
			name:    "tag starting with colon (invalid key)",
			tag:     ":",
			wantErr: errors.New("invalid syntax for struct tag key"),
		},
		{
			name:    "tag with colon not followed by quote (invalid value)",
			tag:     "json:foo",
			wantErr: errors.New("invalid syntax for struct tag value"),
		},
		{
			name:    "tag with unclosed quote (invalid value)",
			tag:     `json:"foo`,
			wantErr: errors.New("invalid syntax for struct tag value"),
		},
		{
			name:    "tag with invalid escape sequence (invalid value)",
			tag:     `json:"\x"`,
			wantErr: errors.New("invalid syntax for struct tag value"),
		},
	}

	for _, ts := range test {
		t.Run(ts.name, func(t *testing.T) {
			tags, err := structtag.Parse(ts.tag)

			if ts.wantErr != nil {
				if err == nil || err.Error() != ts.wantErr.Error() {
					t.Errorf("unexpected error: got = %v, want = %v", err, ts.wantErr)
				}
				return
			}

			if ts.want == nil {
				if tags != nil {
					t.Errorf("expected tags to be nil, but got %#v", tags)
				}
				return
			}

			if !reflect.DeepEqual(ts.want, tags) {
				t.Errorf("got = %#v, want = %#v", tags, ts.want)
			}
		})
	}
}
