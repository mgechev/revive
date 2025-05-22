package test

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

// configureRule configures the given rule with the given configuration
// if the rule implements the ConfigurableRule interface
func configureRule(t *testing.T, rule lint.Rule, arguments lint.Arguments) {
	t.Helper()

	cr, ok := rule.(lint.ConfigurableRule)
	if !ok {
		return
	}

	err := cr.Configure(arguments)
	if err != nil {
		t.Fatalf("Cannot configure rule %s: %v", rule.Name(), err)
	}
}

var testDataDir = filepath.Join("..", "testdata")

func testRule(t *testing.T, filename string, rule lint.Rule, config ...*lint.RuleConfig) {
	t.Helper()

	baseDir := filepath.Join(testDataDir, filepath.Dir(filename))
	filename = filepath.Base(filename) + ".go"
	fullFilePath := filepath.Join(baseDir, filename)
	_, err := os.Stat(fullFilePath)
	if err != nil {
		t.Fatalf("Cannot get file info for %s: %v", rule.Name(), err)
	}
	testRuleOnFiles(t, []string{fullFilePath}, rule, config...)
}

// testRuleOnDir tests the given rule on all the Go files under the given dir
func testRuleOnDir(t *testing.T, dir string, rule lint.Rule, config ...*lint.RuleConfig) {
	pattern := filepath.Join(testDataDir, dir, "...")
	packages, err := revivelib.GetPackages([]string{pattern}, revivelib.ArrayFlags{})
	if err != nil {
		t.Fatalf("Can not retrieve packages under test: %v", err)
	}

	if len(packages) == 0 {
		t.Fatalf("Can not retrieve packages for pattern %v", pattern)
	}

	files := packages[0]
	if len(packages) == 0 {
		t.Fatalf("Can not retrieve files for pattern %v", pattern)
	}

	testRuleOnFiles(t, files, rule, config...)
}

func testRuleOnFiles(t *testing.T, files []string, rule lint.Rule, config ...*lint.RuleConfig) {
	var ruleConfig lint.RuleConfig
	c := map[string]lint.RuleConfig{}
	if len(config) > 0 {
		ruleConfig = *config[0]
		c[rule.Name()] = ruleConfig
	}
	configureRule(t, rule, ruleConfig.Arguments)

	if len(parseInstructions(t, files)) == 0 {
		assertSuccess(t, files, []lint.Rule{rule}, c)
		return
	}
	assertFailures(t, files, []lint.Rule{rule}, c)
}

func assertSuccess(t *testing.T, files []string, rules []lint.Rule, config map[string]lint.RuleConfig) error {
	t.Helper()

	l := lint.New(os.ReadFile, 0)

	ps, err := l.Lint([][]string{files}, rules, lint.Config{
		Rules: config,
	})
	if err != nil {
		return err
	}

	failures := ""
	for p := range ps {
		failures += p.String() + "\n"
	}
	if failures != "" {
		t.Errorf("Expected the rule to pass but got the following failures: %s", failures)
	}
	return nil
}

func assertFailures(t *testing.T, files []string, rules []lint.Rule, config map[string]lint.RuleConfig) error {
	t.Helper()

	l := lint.New(os.ReadFile, 0)

	ins := parseInstructions(t, files)

	ps, err := l.Lint([][]string{files}, rules, lint.Config{
		Rules: config,
	})
	if err != nil {
		return err
	}

	failures := map[string][]lint.Failure{} // file -> failures
	for f := range ps {
		failures[f.GetFilename()] = append(failures[f.GetFilename()], f)
	}

	for file, ins := range ins {
		failures := failures[file]
		for _, in := range ins {
			ok := false
			for i, p := range failures {
				if p.Position.Start.Line != in.Line {
					continue
				}

				if in.Match == p.Failure {
					// check replacement if we are expecting one
					if in.Replacement != "" {
						// ignore any inline comments, since that would be recursive
						r := p.ReplacementLine
						if i := strings.Index(r, " //"); i >= 0 {
							r = r[:i]
						}
						if r != in.Replacement {
							t.Errorf("Lint failed at %s:%d; got replacement %q, want %q", file, in.Line, r, in.Replacement)
						}
					}

					if in.Confidence > 0 {
						if in.Confidence != p.Confidence {
							t.Errorf("Lint failed at %s:%d; got confidence %f, want %f", file, in.Line, p.Confidence, in.Confidence)
						}
					}

					// remove this problem from ps
					copy(failures[i:], failures[i+1:])
					failures = failures[:len(failures)-1]

					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("Lint failed at %s:%d; /%v/ did not match", file, in.Line, in.Match)
			}
		}
		for _, p := range failures {
			t.Errorf("Unexpected problem at %s:%d: %v", file, p.Position.Start.Line, p.Failure)
		}
	}
	return nil
}

type instruction struct {
	Line        int     // the line number this applies to
	Match       string  // which pattern to match
	Replacement string  // what the suggested replacement line should be
	RuleName    string  // what rule we use
	Category    string  // which category
	Confidence  float64 // confidence level
}

// JSONInstruction structure used when we parse json object instead of classic MATCH string
type JSONInstruction struct {
	Match      string  `json:"MATCH"`
	Category   string  `json:"Category"`
	Confidence float64 `json:"Confidence"`
}

// parseInstructions parses instructions from the comments in a Go source file.
// It returns nil if none were parsed.
func parseInstructions(t *testing.T, files []string) map[string][]instruction {
	t.Helper()

	result := map[string][]instruction{}
	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Bad filename path in test: %v", err)
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
		if err != nil {
			t.Fatalf("Test file %v does not parse: %v", file, err)
		}
		var ins []instruction
		for _, cg := range f.Comments {
			ln := fset.Position(cg.Pos()).Line
			raw := cg.Text()
			for _, line := range strings.Split(raw, "\n") {
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "ignore") {
					continue
				}
				if line == "OK" && ins == nil {
					// so our return value will be non-nil
					ins = []instruction{}
					continue
				}
				switch extractDataMode(line) {
				case "json":
					jsonInst, err := extractInstructionFromJSON(strings.TrimPrefix(line, "json:"), ln)
					if err != nil {
						t.Fatalf("At %v:%d: %v", file, ln, err)
					}
					ins = append(ins, jsonInst)
				case "classic":
					match, err := extractPattern(line)
					if err != nil {
						t.Fatalf("At %v:%d: %v", file, ln, err)
					}
					matchLine := ln
					if i := strings.Index(line, "MATCH:"); i >= 0 {
						// This is a match for a different line.
						lns := strings.TrimPrefix(line[i:], "MATCH:")
						lns = lns[:strings.Index(lns, " ")]
						matchLine, err = strconv.Atoi(lns)
						if err != nil {
							t.Fatalf("Bad match line number %q at %v:%d: %v", lns, file, ln, err)
						}
					}
					var repl string
					if r, ok := extractReplacement(line); ok {
						repl = r
					}
					ins = append(ins, instruction{
						Line:        matchLine,
						Match:       match,
						Replacement: repl,
					})
				}
			}
		}
		if len(ins) != 0 {
			result[file] = ins
		}
	}

	return result
}

func extractInstructionFromJSON(line string, lineNumber int) (instruction, error) {
	// Use the json.Unmarshal function to parse the JSON into the struct
	var jsonInst JSONInstruction
	if err := json.Unmarshal([]byte(line), &jsonInst); err != nil {
		return instruction{}, fmt.Errorf("parsing json instruction: %w", err)
	}

	ins := instruction{
		Match:      jsonInst.Match,
		Confidence: jsonInst.Confidence,
		Category:   jsonInst.Category,
		Line:       lineNumber,
	}
	return ins, nil
}

func extractDataMode(line string) string {
	if strings.HasPrefix(line, "json") {
		return "json"
	}
	if strings.Contains(line, "MATCH") {
		return "classic"
	}

	return ""
}

func extractPattern(line string) (string, error) {
	a, b := strings.Index(line, "/"), strings.LastIndex(line, "/")
	if a == -1 || a == b {
		return "", fmt.Errorf("malformed match instruction %q", line)
	}
	return line[a+1 : b], nil
}

func extractReplacement(line string) (string, bool) {
	// Look for this:  / -> `
	// (the end of a match and start of a backtick string),
	// and then the closing backtick.
	const start = "/ -> `"
	a, b := strings.Index(line, start), strings.LastIndex(line, "`")
	if a < 0 || a > b {
		return "", false
	}
	return line[a+len(start) : b], true
}

func srcLine(src []byte, p token.Position) string {
	// Run to end of line in both directions if not at line start/end.
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}

// TestLine tests srcLine function
func TestLine(t *testing.T) {
	tests := []struct {
		src    string
		offset int
		want   string
	}{
		{"single line file", 5, "single line file"},
		{"single line file with newline\n", 5, "single line file with newline\n"},
		{"first\nsecond\nthird\n", 2, "first\n"},
		{"first\nsecond\nthird\n", 9, "second\n"},
		{"first\nsecond\nthird\n", 14, "third\n"},
		{"first\nsecond\nthird with no newline", 16, "third with no newline"},
		{"first byte\n", 0, "first byte\n"},
	}
	for _, test := range tests {
		got := srcLine([]byte(test.src), token.Position{Offset: test.offset})
		if got != test.want {
			t.Errorf("srcLine(%q, offset=%d) = %q, want %q", test.src, test.offset, got, test.want)
		}
	}
}

// exportedType reports whether typ is an exported type.
// It is imprecise, and will err on the side of returning true,
// such as for composite types.
func exportedType(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Named:
		// Builtin types have no package.
		return t.Obj().Pkg() == nil || t.Obj().Exported()
	case *types.Map:
		return exportedType(t.Key()) && exportedType(t.Elem())
	case interface {
		Elem() types.Type
	}: // array, slice, pointer, chan
		return exportedType(t.Elem())
	}
	// Be conservative about other types, such as struct, interface, etc.
	return true
}

// TestExportedType tests exportedType function
func TestExportedType(t *testing.T) {
	tests := []struct {
		typString string
		exp       bool
	}{
		{"int", true},
		{"string", false}, // references the shadowed builtin "string"
		{"T", true},
		{"t", false},
		{"*T", true},
		{"*t", false},
		{"map[int]complex128", true},
	}
	for _, test := range tests {
		src := `package foo; type T int; type t int; type string struct{}`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "foo.go", src, 0)
		if err != nil {
			t.Fatalf("Parsing %q: %v", src, err)
		}
		// use the package name as package path
		config := &types.Config{}
		pkg, err := config.Check(file.Name.Name, fset, []*ast.File{file}, nil)
		if err != nil {
			t.Fatalf("Type checking %q: %v", src, err)
		}
		tv, err := types.Eval(fset, pkg, token.NoPos, test.typString)
		if err != nil {
			t.Errorf("types.Eval(%q): %v", test.typString, err)
			continue
		}
		if got := exportedType(tv.Type); got != test.exp {
			t.Errorf("exportedType(%v) = %t, want %t", tv.Type, got, test.exp)
		}
	}
}
