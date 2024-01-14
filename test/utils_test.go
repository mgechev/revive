package test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/pkg/errors"
)

func testRule(t *testing.T, filename string, rule lint.Rule, config ...*lint.RuleConfig) {
	baseDir := "../testdata/"
	filename = filename + ".go"
	src, err := os.ReadFile(baseDir + filename)
	if err != nil {
		t.Fatalf("Bad filename path in test for %s: %v", rule.Name(), err)
	}
	stat, err := os.Stat(baseDir + filename)
	if err != nil {
		t.Fatalf("Cannot get file info for %s: %v", rule.Name(), err)
	}
	c := map[string]lint.RuleConfig{}
	if config != nil {
		c[rule.Name()] = *config[0]
	}
	if parseInstructions(t, filename, src) == nil {
		assertSuccess(t, baseDir, stat, []lint.Rule{rule}, c)
		return
	}
	assertFailures(t, baseDir, stat, src, []lint.Rule{rule}, c)
}

func assertSuccess(t *testing.T, baseDir string, fi os.FileInfo, rules []lint.Rule, config map[string]lint.RuleConfig) error {
	l := lint.New(func(file string) ([]byte, error) {
		return os.ReadFile(baseDir + file)
	}, 0)

	ps, err := l.Lint([][]string{{fi.Name()}}, rules, lint.Config{
		Rules: config,
	})
	if err != nil {
		return err
	}

	failures := ""
	for p := range ps {
		failures += p.Failure
	}
	if failures != "" {
		t.Errorf("Expected the rule to pass but got the following failures: %s", failures)
	}
	return nil
}

func assertFailures(t *testing.T, baseDir string, fi os.FileInfo, src []byte, rules []lint.Rule, config map[string]lint.RuleConfig) error {
	l := lint.New(func(file string) ([]byte, error) {
		return os.ReadFile(baseDir + file)
	}, 0)

	ins := parseInstructions(t, fi.Name(), src)
	if ins == nil {
		return errors.Errorf("Test file %v does not have instructions", fi.Name())
	}

	ps, err := l.Lint([][]string{{fi.Name()}}, rules, lint.Config{
		Rules: config,
	})
	if err != nil {
		return err
	}

	failures := []lint.Failure{}
	for f := range ps {
		failures = append(failures, f)
	}

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
						t.Errorf("Lint failed at %s:%d; got replacement %q, want %q", fi.Name(), in.Line, r, in.Replacement)
					}
				}

				if in.Confidence > 0 {

					if in.Confidence != p.Confidence {
						t.Errorf("Lint failed at %s:%d; got confidence %f, want %f", fi.Name(), in.Line, p.Confidence, in.Confidence)
					}

				}

				// remove this problem from ps
				copy(failures[i:], failures[i+1:])
				failures = failures[:len(failures)-1]

				// t.Logf("/%v/ matched at %s:%d", in.Match, fi.Name(), in.Line)
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("Lint failed at %s:%d; /%v/ did not match", fi.Name(), in.Line, in.Match)
		}
	}
	for _, p := range failures {
		t.Errorf("Unexpected problem at %s:%d: %v", fi.Name(), p.Position.Start.Line, p.Failure)
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

// JSONInstruction structure used when we parse json object insted of classic MATCH string
type JSONInstruction struct {
	Match      string  `json:"MATCH"`
	Category   string  `json:"Category"`
	Confidence float64 `json:"Confidence"`
}

// parseInstructions parses instructions from the comments in a Go source file.
// It returns nil if none were parsed.
func parseInstructions(t *testing.T, filename string, src []byte) []instruction {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Test file %v does not parse: %v", filename, err)
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
				ins = make([]instruction, 0)
				continue
			}
			switch extractDataMode(line) {
			case "json":
				jsonInst, err := extractInstructionFromJSON(strings.TrimPrefix(line, "json:"), ln)
				if err != nil {
					t.Fatalf("At %v:%d: %v", filename, ln, err)
				}
				ins = append(ins, jsonInst)
				break
			case "classic":
				match, err := extractPattern(line)
				if err != nil {
					t.Fatalf("At %v:%d: %v", filename, ln, err)
				}
				matchLine := ln
				if i := strings.Index(line, "MATCH:"); i >= 0 {
					// This is a match for a different line.
					lns := strings.TrimPrefix(line[i:], "MATCH:")
					lns = lns[:strings.Index(lns, " ")]
					matchLine, err = strconv.Atoi(lns)
					if err != nil {
						t.Fatalf("Bad match line number %q at %v:%d: %v", lns, filename, ln, err)
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
				break
			}

		}
	}
	return ins
}

func extractInstructionFromJSON(line string, lineNumber int) (instruction, error) {
	// Use the json.Unmarshal function to parse the JSON into the struct
	var jsonInst JSONInstruction
	if err := json.Unmarshal([]byte(line), &jsonInst); err != nil {
		fmt.Println("Error parsing JSON:", err)
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
func TestLine(t *testing.T) { //revive:disable-line:exported
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

// TestLintName tests lint.Name function
func TestLintName(t *testing.T) { //revive:disable-line:exported
	tests := []struct {
		name, want string
	}{
		{"foo_bar", "fooBar"},
		{"foo_bar_baz", "fooBarBaz"},
		{"Foo_bar", "FooBar"},
		{"foo_WiFi", "fooWiFi"},
		{"id", "id"},
		{"Id", "ID"},
		{"foo_id", "fooID"},
		{"fooId", "fooID"},
		{"fooUid", "fooUID"},
		{"idFoo", "idFoo"},
		{"uidFoo", "uidFoo"},
		{"midIdDle", "midIDDle"},
		{"APIProxy", "APIProxy"},
		{"ApiProxy", "APIProxy"},
		{"apiProxy", "apiProxy"},
		{"_Leading", "_Leading"},
		{"___Leading", "_Leading"},
		{"trailing_", "trailing"},
		{"trailing___", "trailing"},
		{"a_b", "aB"},
		{"a__b", "aB"},
		{"a___b", "aB"},
		{"Rpc1150", "RPC1150"},
		{"case3_1", "case3_1"},
		{"case3__1", "case3_1"},
		{"IEEE802_16bit", "IEEE802_16bit"},
		{"IEEE802_16Bit", "IEEE802_16Bit"},
	}
	for _, test := range tests {
		got := lint.Name(test.name, nil, nil)
		if got != test.want {
			t.Errorf("lintName(%q) = %q, want %q", test.name, got, test.want)
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
func TestExportedType(t *testing.T) { //revive:disable-line:exported
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

var (
	genHdr = []byte("// Code generated ")
	genFtr = []byte(" DO NOT EDIT.")
)

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.HasSuffix(b, genFtr) && len(b) >= len(genHdr)+len(genFtr) {
			return true
		}
	}
	return false
}

// TestIsGenerated tests isGenerated function
func TestIsGenerated(t *testing.T) { //revive:disable-line:exported
	tests := []struct {
		source    string
		generated bool
	}{
		{"// Code Generated by some tool. DO NOT EDIT.", false},
		{"// Code generated by some tool. DO NOT EDIT.", true},
		{"// Code generated by some tool. DO NOT EDIT", false},
		{"// Code generated  DO NOT EDIT.", true},
		{"// Code generated DO NOT EDIT.", false},
		{"\t\t// Code generated by some tool. DO NOT EDIT.\npackage foo\n", false},
		{"// Code generated by some tool. DO NOT EDIT.\npackage foo\n", true},
		{"package foo\n// Code generated by some tool. DO NOT EDIT.\ntype foo int\n", true},
		{"package foo\n // Code generated by some tool. DO NOT EDIT.\ntype foo int\n", false},
		{"package foo\n// Code generated by some tool. DO NOT EDIT. \ntype foo int\n", false},
		{"package foo\ntype foo int\n// Code generated by some tool. DO NOT EDIT.\n", true},
		{"package foo\ntype foo int\n// Code generated by some tool. DO NOT EDIT.", true},
	}

	for i, test := range tests {
		got := isGenerated([]byte(test.source))
		if got != test.generated {
			t.Errorf("test %d, isGenerated() = %v, want %v", i, got, test.generated)
		}
	}
}
