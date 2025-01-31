package lint

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestFile_disabledIntervals(t *testing.T) {
	buildCommentGroups := func(comments ...string) []*ast.CommentGroup {
		commentGroups := make([]*ast.CommentGroup, 0, len(comments))
		for _, c := range comments {
			commentGroups = append(commentGroups, &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: c},
				},
			})
		}
		return commentGroups
	}

	tests := []struct {
		name     string
		comments []*ast.CommentGroup
		expected disabledIntervalsMap
	}{
		{
			name:     "no directives",
			comments: buildCommentGroups("// some comment"),
			expected: disabledIntervalsMap{},
		},
		{
			name:     "disable rule",
			comments: buildCommentGroups("//revive:disable:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
						},
						To: token.Position{
							Filename: "test.go",
							Line:     2147483647,
						},
					},
				},
			},
		},
		{
			name:     "enable rule",
			comments: buildCommentGroups("//revive:enable:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {},
			},
		},
		{
			name:     "disable and enable rule",
			comments: buildCommentGroups("//revive:disable:rule1", "//revive:enable:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
						},
						To: token.Position{
							Filename: "test.go",
						},
					},
				},
			},
		},
		{
			name:     "disable-line rule",
			comments: buildCommentGroups("//revive:disable-line:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
						},
						To: token.Position{
							Filename: "test.go",
						},
					},
				},
			},
		},
		{
			name:     "enable-line rule",
			comments: buildCommentGroups("//revive:enable-line:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
						},
						To: token.Position{
							Filename: "test.go",
							Line:     2147483647,
						},
					},
				},
			},
		},
		{
			name:     "disable-next-line rule",
			comments: buildCommentGroups("//revive:disable-next-line:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
							Line:     1,
						},
						To: token.Position{
							Filename: "test.go",
							Line:     1,
						},
					},
				},
			},
		},
		{
			name:     "enable-next-line rule",
			comments: buildCommentGroups("//revive:enable-next-line:rule1"),
			expected: disabledIntervalsMap{
				"rule1": {
					{
						RuleName: "rule1",
						From: token.Position{
							Filename: "test.go",
							Line:     1,
						},
						To: token.Position{
							Filename: "test.go",
							Line:     2147483647,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Name: "test.go",
				Pkg: &Package{
					fset: token.NewFileSet(),
				},
				AST: &ast.File{
					Comments: tt.comments,
				},
			}
			got := f.disabledIntervals(nil, false, make(chan Failure, 10))
			if len(got) != len(tt.expected) {
				t.Errorf("disabledIntervals() = got %v, want %v", got, tt.expected)
			}
			for rule, intervals := range got {
				expectedIntervals, ok := tt.expected[rule]
				if !ok {
					t.Errorf("unexpected rule %q", rule)
					continue
				}
				if len(intervals) != len(expectedIntervals) {
					t.Errorf("intervals for rule %q = got %+v, want %+v", rule, intervals, expectedIntervals)
					continue
				}
				for i, interval := range intervals {
					if interval != expectedIntervals[i] {
						t.Errorf("interval %d for rule %q = got %+v, want %+v", i, rule, interval, expectedIntervals[i])
					}
				}
			}
		})
	}
}
