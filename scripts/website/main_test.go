package main

import (
	"strings"
	"testing"
)

func TestSplitSections(t *testing.T) {
	doc := `# Description of available rules

List of all available rules.

<!-- toc -->

- [add-constant](#add-constant)

<!-- tocstop -->

## Configuration options format

Options use ` + "`kebab-case`" + `.

## add-constant

_Go version_: 1.0.

_Description_: Suggests using constants.

` + "```go" + `
// ## looks-like-a-heading but is inside a fence
` + "```" + `

## bare-return

_Description_: Warns on bare returns.
`
	sections := splitSections(doc)
	if len(sections) != 3 {
		t.Fatalf("got %d sections, want 3: %+v", len(sections), sections)
	}

	wantNames := []string{"configuration-options-format", "add-constant", "bare-return"}
	for i, want := range wantNames {
		if sections[i].name != want {
			t.Errorf("section %d name = %q, want %q", i, sections[i].name, want)
		}
	}
	if got := sections[0].heading; got != "Configuration options format" {
		t.Errorf("section 0 heading = %q, want %q", got, "Configuration options format")
	}

	if body := sections[0].body; body != "Options use `kebab-case`.\n" {
		t.Errorf("unexpected configuration-options-format body: %q", body)
	}
	if body := sections[1].body; !strings.Contains(body, "looks-like-a-heading") {
		t.Errorf("fenced pseudo-heading split the add-constant section: %q", body)
	}
	if body := sections[1].body; strings.Contains(body, "## add-constant") {
		t.Errorf("section body must not contain its own heading: %q", body)
	}
}

func TestSplitSectionsEmptyDoc(t *testing.T) {
	if sections := splitSections("no headings at all\n"); sections != nil {
		t.Errorf("got %+v, want nil", sections)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Configuration options format", "configuration-options-format"},
		{"add-constant", "add-constant"},
		{"  spaced  ", "spaced"},
	}
	for _, tt := range tests {
		if got := slugify(tt.in); got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRewriteRuleLinks(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{
			name: "fragment link to another rule stays a same-page anchor",
			in:   "handled by [indent-error-flow](#indent-error-flow).",
			want: "handled by [indent-error-flow](#indent-error-flow).",
		},
		{
			name: "explicit RULES_DESCRIPTIONS.md link becomes a same-page anchor",
			in:   "see [var-naming](./RULES_DESCRIPTIONS.md#var-naming)",
			want: "see [var-naming](#var-naming)",
		},
		{
			name: "README link with fragment",
			in:   "see the [README](./README.md#configuration)",
			want: "see the [README](/docs/)",
		},
		{
			name: "absolute URL untouched",
			in:   "see [the blog](https://go.dev/blog/package-names#section) and [wiki](https://en.wikipedia.org/wiki/X_(y))",
			want: "see [the blog](https://go.dev/blog/package-names#section) and [wiki](https://en.wikipedia.org/wiki/X_(y))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rewriteRuleLinks(tt.in); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTransformReadme(t *testing.T) {
	doc := `# revive

[![Build Status](https://example.com/badge.svg)](https://example.com/build)

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
</p>

<!-- toc -->

- [Usage](#usage)

<!-- tocstop -->

## Usage

List of all [available rules](./RULES_DESCRIPTIONS.md).

| [` + "`add-constant`" + `](./RULES_DESCRIPTIONS.md#add-constant) | map |

See [DEVELOPING.md](./DEVELOPING.md) and [CONTRIBUTING.md](CONTRIBUTING.md).
See [Configuration](#configuration) and ![demo](assets/demo.svg).
![friendly](/assets/formatter-friendly.png)
Look at [this file](/formatter/json.go) for an example.
`
	got := transformReadme(doc)

	for _, unwanted := range []string{
		"# revive\n",
		"<!-- toc -->",
		"<!-- tocstop -->",
		"- [Usage](#usage)",
		"assets/logo.png",
		"./RULES_DESCRIPTIONS.md",
	} {
		if strings.Contains(got, unwanted) {
			t.Errorf("output still contains %q:\n%s", unwanted, got)
		}
	}

	for _, wanted := range []string{
		"[![Build Status](https://example.com/badge.svg)](https://example.com/build)",
		`<img src="/images/logo.png"`,
		"[available rules](https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md)",
		"](/r/#add-constant)",
		"[DEVELOPING.md](https://github.com/mgechev/revive/blob/master/DEVELOPING.md)",
		"[CONTRIBUTING.md](https://github.com/mgechev/revive/blob/master/CONTRIBUTING.md)",
		"[Configuration](#configuration)",
		"![demo](/images/demo.svg)",
		"![friendly](/images/formatter-friendly.png)",
		"[this file](https://github.com/mgechev/revive/blob/master/formatter/json.go)",
		"## Usage",
	} {
		if !strings.Contains(got, wanted) {
			t.Errorf("output does not contain %q:\n%s", wanted, got)
		}
	}

	if strings.HasPrefix(got, "\n") {
		t.Errorf("output starts with a blank line:\n%s", got)
	}
}

func TestPageContent(t *testing.T) {
	got := pageContent(frontMatter{title: "add-constant", description: `Suggests "constants".`}, "Body.")
	want := "---\n" +
		"title: \"add-constant\"\n" +
		"description: \"Suggests \\\"constants\\\".\"\n" +
		"---\n\n" +
		"Body.\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}

	got = pageContent(frontMatter{title: "Rule descriptions", pageType: "docs"}, "Body.\n")
	want = "---\n" +
		"title: \"Rule descriptions\"\n" +
		"type: \"docs\"\n" +
		"---\n\n" +
		"Body.\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCompareRuleSets(t *testing.T) {
	if err := compareRuleSets([]string{"a", "b"}, []string{"b", "a"}); err != nil {
		t.Errorf("equal sets: got error %v, want nil", err)
	}

	err := compareRuleSets([]string{"documented-only", "both"}, []string{"registered-only", "both"})
	if err == nil {
		t.Fatal("diverged sets: got nil, want error")
	}
	if got := err.Error(); !strings.Contains(got, "registered but not documented: registered-only") ||
		!strings.Contains(got, "documented but not registered: documented-only") {
		t.Errorf("error does not list both directions: %v", err)
	}
}

func TestRulesPageContent(t *testing.T) {
	sections := []section{
		{name: "configuration-options-format", heading: "Configuration options format", body: "Options use `kebab-case`.\n"},
		{name: "add-constant", heading: "add-constant", body: "_Description_: See [var-naming](./RULES_DESCRIPTIONS.md#var-naming).\n"},
		{name: "bare-return", heading: "bare-return", body: "_Description_: Warns on bare returns.\n"},
	}
	got := rulesPageContent(sections)

	for _, wanted := range []string{
		"## Configuration options format\n",
		"## add-constant\n",
		"## bare-return\n",
		"[var-naming](#var-naming)",
	} {
		if !strings.Contains(got, wanted) {
			t.Errorf("page body does not contain %q:\n%s", wanted, got)
		}
	}

	addConstant := strings.Index(got, "## add-constant")
	bareReturn := strings.Index(got, "## bare-return")
	if addConstant < 0 || bareReturn < 0 || addConstant > bareReturn {
		t.Errorf("sections are not in documentation order:\n%s", got)
	}
}
