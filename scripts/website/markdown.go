package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// section is a "## " section of RULES_DESCRIPTIONS.md.
type section struct {
	name    string // slugified heading, e.g. "add-constant"
	heading string // raw heading text, e.g. "Configuration options format"
	body    string // section content without the heading line
}

// splitSections drops everything before the first "## " heading (title, intro,
// table of contents) and splits the rest of the document into sections, one
// per "## " heading. Headings inside code fences are ignored.
func splitSections(doc string) []section {
	var (
		sections []section
		name     string
		heading  string
		buf      []string
		inFence  bool
		started  bool
	)
	flush := func() {
		if !started {
			return
		}
		body := strings.Trim(strings.Join(buf, "\n"), "\n")
		if body != "" {
			body += "\n"
		}
		sections = append(sections, section{name: name, heading: heading, body: body})
	}
	for line := range strings.SplitSeq(doc, "\n") {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(line, "## ") {
			flush()
			heading = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			name = slugify(heading)
			started = true
			buf = buf[:0]
			continue
		}
		if started {
			buf = append(buf, line)
		}
	}
	flush()
	return sections
}

// slugify converts a section heading to its GitHub anchor/slug form,
// e.g. "Configuration options format" -> "configuration-options-format".
func slugify(heading string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(heading)), " ", "-")
}

var (
	// rulesFileLinkRE matches "./RULES_DESCRIPTIONS.md#<rule-name>" cross-references.
	// Plain fragment links like "](#other-rule)" need no rewrite: all rules share the /r/ page.
	rulesFileLinkRE = regexp.MustCompile(`\]\(\./RULES_DESCRIPTIONS\.md#([A-Za-z0-9_-]+)\)`)
	// readmeLinkRE matches "./README.md" links, with or without a fragment.
	readmeLinkRE = regexp.MustCompile(`\]\(\./README\.md[^)]*\)`)
	// relativeMdLinkRE matches repository-relative Markdown file links in README.md.
	relativeMdLinkRE = regexp.MustCompile(`\]\((?:\./)?([A-Za-z0-9][A-Za-z0-9_./-]*\.md)\)`)
	// repoFileLinkRE matches root-relative Go source file links in README.md.
	repoFileLinkRE = regexp.MustCompile(`\]\(/([A-Za-z0-9][A-Za-z0-9_./-]*\.go)\)`)
	// assetRefRE matches assets/ image references in README.md links and img tags.
	assetRefRE = regexp.MustCompile(`(\]\(|src=")(?:\.?/)?assets/`)
)

// rewriteRuleLinks rewrites the cross-references of a RULES_DESCRIPTIONS.md
// section to the URL scheme of the website. Absolute URLs and same-page
// fragment links are left untouched.
func rewriteRuleLinks(body string) string {
	body = rulesFileLinkRE.ReplaceAllString(body, "](#${1})")
	body = readmeLinkRE.ReplaceAllString(body, "](/docs/)")
	return body
}

// transformReadme turns the README.md content into the body of the /docs/
// page: it drops the leading H1 and the table-of-contents block, and rewrites
// repository-relative references. Absolute URLs are left untouched.
func transformReadme(doc string) string {
	doc = stripTocBlock(doc)
	doc = stripLeadingH1(doc)
	doc = rewriteReadmeLinks(doc)
	return strings.TrimLeft(doc, "\n")
}

// stripTocBlock removes the markdown-toc block, "<!-- toc -->" through
// "<!-- tocstop -->" inclusive. The document is returned unchanged if the
// markers are absent.
func stripTocBlock(doc string) string {
	var kept []string
	inToc := false
	for line := range strings.SplitSeq(doc, "\n") {
		switch {
		case strings.TrimSpace(line) == "<!-- toc -->":
			inToc = true
		case strings.TrimSpace(line) == "<!-- tocstop -->":
			inToc = false
		case !inToc:
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// stripLeadingH1 removes the document's H1 heading if it is the first
// non-blank line.
func stripLeadingH1(doc string) string {
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			return strings.Join(append(lines[:i:i], lines[i+1:]...), "\n")
		}
		break
	}
	return doc
}

func rewriteReadmeLinks(doc string) string {
	doc = rulesFileLinkRE.ReplaceAllString(doc, "](/r/#${1})")
	doc = relativeMdLinkRE.ReplaceAllString(doc, "](https://github.com/mgechev/revive/blob/master/${1})")
	doc = repoFileLinkRE.ReplaceAllString(doc, "](https://github.com/mgechev/revive/blob/master/${1})")
	doc = assetRefRE.ReplaceAllString(doc, "${1}/images/")
	return doc
}

// frontMatter is the YAML front matter of a generated page.
type frontMatter struct {
	title       string
	description string
	pageType    string // Hugo page type, e.g. "docs" for Hextra's sidebar layout
}

// pageContent renders a Hugo content file: YAML front matter followed by body.
func pageContent(fm frontMatter, body string) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", strconv.Quote(fm.title))
	if fm.description != "" {
		fmt.Fprintf(&b, "description: %s\n", strconv.Quote(fm.description))
	}
	if fm.pageType != "" {
		fmt.Fprintf(&b, "type: %s\n", strconv.Quote(fm.pageType))
	}
	b.WriteString("---\n\n")
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteString("\n")
	}
	return b.String()
}

// rulesPageContent renders the body of the single /r/ page: a short intro
// followed by every RULES_DESCRIPTIONS.md section as an "## " heading in
// documentation order. Heading anchors are GitHub-style slugs, so legacy
// /r#<rule-name> and /r/#<rule-name> links keep working with no redirects.
func rulesPageContent(sections []section) string {
	var b strings.Builder
	b.WriteString("This is the complete list of the linting rules provided by `revive`,\n")
	b.WriteString("with their behavior, configuration options, and examples.\n")
	for _, sec := range sections {
		fmt.Fprintf(&b, "\n## %s\n\n", sec.heading)
		b.WriteString(rewriteRuleLinks(sec.body))
	}
	return b.String()
}
