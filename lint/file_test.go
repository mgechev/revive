package lint

import (
	"bytes"
	"go/ast"
	"go/token"
	"log/slog"
	"slices"
	"strings"
	"testing"
)

type fakeRule struct {
	name     string
	failures []Failure
}

var _ Rule = (*fakeRule)(nil)

func (r *fakeRule) Name() string { return r.name }

func (r *fakeRule) Apply(*File, Arguments) []Failure {
	return slices.Clone(r.failures)
}

func TestFile_lint_internalFailureDoesNotAbortOtherRules(t *testing.T) {
	rules := []Rule{
		&fakeRule{
			name: "typecheck-internal-failure-rule",
			failures: []Failure{
				NewInternalFailure("simulated type-check failure"),
			},
		},
		&fakeRule{
			name: "normal-rule",
			failures: []Failure{
				{
					Confidence: 1,
					Failure:    "must reach the channel",
				},
			},
		},
	}

	cfg := Config{
		Confidence: 0.8,
		Rules: RulesConfig{
			"typecheck-internal-failure-rule": {},
			"normal-rule":                     {},
		},
	}

	var logBuf bytes.Buffer
	f := &File{
		Name:   "test.go",
		Pkg:    &Package{fset: token.NewFileSet()},
		AST:    &ast.File{},
		logger: slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}

	failures := make(chan Failure, 4)
	if err := f.lint(rules, cfg, failures); err != nil {
		t.Fatal("unexpected error from linting:", err)
	}
	close(failures)

	var got []Failure
	for failure := range failures {
		got = append(got, failure)
	}

	if len(got) != 1 {
		t.Fatalf("expected exactly 1 failure to be reported, got %d: %+v", len(got), got)
	}
	if got[0].RuleName != "normal-rule" {
		t.Errorf("expected failure from %q, got %q", "normal-rule", got[0].RuleName)
	}

	logged := logBuf.String()
	for _, want := range []string{
		"level=WARN",
		`msg="rule skipped due to internal failure"`,
		"rule=typecheck-internal-failure-rule",
		"file=test.go",
		`failure="simulated type-check failure"`,
	} {
		if !strings.Contains(logged, want) {
			t.Errorf("expected log output to contain %q, got %q", want, logged)
		}
	}
}

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
				logger: slog.New(slog.DiscardHandler),
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
