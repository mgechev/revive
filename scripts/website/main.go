// Program website generates the content of the revive website (https://revive.run/)
// from the documentation that lives in this repository.
//
// Run it from the repository root:
//
//	go run ./scripts/website
//
// It reads README.md and RULES_DESCRIPTIONS.md, cross-checks the documented
// rules against the rules registered in the config package, writes Hugo
// content files under docs/content/, and copies the files from assets/ to
// docs/static/images/.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mgechev/revive/config"
)

// configOptionsSection is the slug of the only section of RULES_DESCRIPTIONS.md
// that does not describe a rule.
const configOptionsSection = "configuration-options-format"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	rulesDoc, err := readTextFile("RULES_DESCRIPTIONS.md")
	if err != nil {
		return err
	}

	sections := splitSections(rulesDoc)
	if len(sections) == 0 {
		return errors.New("no \"## \" sections found in RULES_DESCRIPTIONS.md")
	}

	var ruleSections []section
	hasConfigSection := false
	seen := make(map[string]bool, len(sections))
	for _, sec := range sections {
		if seen[sec.name] {
			return fmt.Errorf("duplicate section %q in RULES_DESCRIPTIONS.md", sec.name)
		}
		seen[sec.name] = true
		if sec.name == configOptionsSection {
			hasConfigSection = true
			continue
		}
		ruleSections = append(ruleSections, sec)
	}
	if !hasConfigSection {
		return fmt.Errorf("section %q not found in RULES_DESCRIPTIONS.md", configOptionsSection)
	}

	if err := checkDrift(ruleSections); err != nil {
		return err
	}

	rDir := filepath.Join("docs", "content", "r")
	if err := resetDir(rDir); err != nil {
		return err
	}

	// All rules live on the single /r/ page, mirroring RULES_DESCRIPTIONS.md:
	// its "## " anchors keep legacy /r#<rule-name> links working natively.
	page := frontMatter{title: "Rule descriptions", description: "List of all available revive rules.", pageType: "docs"}
	if err := writePage(filepath.Join(rDir, "_index.md"), page, rulesPageContent(sections)); err != nil {
		return err
	}

	readme, err := readTextFile("README.md")
	if err != nil {
		return err
	}
	docsDir := filepath.Join("docs", "content", "docs")
	if err := resetDir(docsDir); err != nil {
		return err
	}
	if err := writePage(filepath.Join(docsDir, "_index.md"), frontMatter{title: "Documentation"}, transformReadme(readme)); err != nil {
		return err
	}

	copied, err := copyAssets("assets", filepath.Join("docs", "static", "images"))
	if err != nil {
		return err
	}

	fmt.Printf("Generated r/_index.md (%d rules) and docs/_index.md; copied %d assets.\n",
		len(ruleSections), copied)
	return nil
}

// checkDrift verifies that the rules documented in RULES_DESCRIPTIONS.md and
// the rules registered in the config package are exactly the same set.
func checkDrift(ruleSections []section) error {
	documented := make([]string, len(ruleSections))
	for i, sec := range ruleSections {
		documented[i] = sec.name
	}

	registered := make([]string, 0, len(documented))
	for _, r := range config.AllRules() {
		registered = append(registered, r.Name())
	}

	return compareRuleSets(documented, registered)
}

// compareRuleSets returns an error listing the rules present in only one of
// the two sets, or nil if the sets are equal.
func compareRuleSets(documented, registered []string) error {
	documentedSet := make(map[string]bool, len(documented))
	for _, name := range documented {
		documentedSet[name] = true
	}
	registeredSet := make(map[string]bool, len(registered))
	for _, name := range registered {
		registeredSet[name] = true
	}

	var undocumented, unregistered []string
	for name := range registeredSet {
		if !documentedSet[name] {
			undocumented = append(undocumented, name)
		}
	}
	for name := range documentedSet {
		if !registeredSet[name] {
			unregistered = append(unregistered, name)
		}
	}
	if len(undocumented) == 0 && len(unregistered) == 0 {
		return nil
	}

	slices.Sort(undocumented)
	slices.Sort(unregistered)
	var b strings.Builder
	b.WriteString("RULES_DESCRIPTIONS.md is out of sync with the rules registered in config.AllRules():")
	if len(undocumented) > 0 {
		fmt.Fprintf(&b, "\n  registered but not documented: %s", strings.Join(undocumented, ", "))
	}
	if len(unregistered) > 0 {
		fmt.Fprintf(&b, "\n  documented but not registered: %s", strings.Join(unregistered, ", "))
	}
	return errors.New(b.String())
}

// resetDir recreates dir empty, so reruns never leave stale generated files
// behind (e.g. the page of a renamed rule).
func resetDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("removing %s: %w", dir, err)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	return nil
}

// readTextFile reads a file and normalizes its line endings to "\n".
func readTextFile(path string) (string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // ignore G304: potential file inclusion via variable
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", path, err)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

// writePage writes a Hugo content file: YAML front matter followed by body.
func writePage(path string, fm frontMatter, body string) error {
	if err := os.WriteFile(path, []byte(pageContent(fm, body)), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// copyAssets copies every regular file from srcDir to dstDir and returns the
// number of files copied.
func copyAssets(srcDir, dstDir string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, fmt.Errorf("reading %s: %w", srcDir, err)
	}
	if err := resetDir(dstDir); err != nil {
		return 0, err
	}

	copied := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, entry.Name())) //nolint:gosec // ignore G304: potential file inclusion via variable
		if err != nil {
			return copied, fmt.Errorf("reading asset %s: %w", entry.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(dstDir, entry.Name()), data, 0o600); err != nil { //nolint:gosec // ignore G703: asset names come from the repository's assets/ directory
			return copied, fmt.Errorf("copying asset %s: %w", entry.Name(), err)
		}
		copied++
	}
	return copied, nil
}
