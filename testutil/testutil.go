package testutil

import (
	"fmt"
	"go/token"
	"regexp"
	"testing"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

var anyRe = regexp.MustCompile(`\[\/?@(\w+)\]`)
var closingRe = regexp.MustCompile(`\[\/@(\w+)\]`)

type pos struct {
	start int
	end   int
}

func extractFailures(code string) map[string]*pos {
	res := anyRe.FindAllStringSubmatchIndex(code, -1)
	if len(res) == 0 {
		return nil
	}

	reduce := 0
	started := map[string]*pos{}
	if len(res)%2 != 0 {
		panic("incorrect test annotations")
	}

	for _, el := range res {
		substr := code[el[0]:el[1]]
		name := code[el[2]:el[3]]
		isEnd := closingRe.MatchString(substr)
		if isEnd && started[name] == nil {
			panic("incorrect test annotation; closed before opened: " + name)
		}
		if !isEnd {
			started[name] = &pos{start: el[0] - reduce}
		} else {
			started[name].end = el[0] - reduce
		}
		reduce += el[1] - el[0]
	}

	return started
}

func stripAnnotations(code string) string {
	return anyRe.ReplaceAllString(code, "")
}

// AssertSuccess checks if given rule runs correctly with no failures.
func AssertSuccessWithName(t *testing.T, code, name string, testingRule rule.Rule, args rule.Arguments) {
	annotations := extractFailures(code)
	if annotations != nil {
		t.Errorf("There should be no failure annotations when verifying successful rule analysis")
	}

	var fileSet token.FileSet
	file, err := file.New(name, []byte(stripAnnotations(code)), &fileSet)
	if err != nil {
		t.Errorf("Cannot parse testing file: %s", err.Error())
	}
	failures := testingRule.Apply(file, args)
	failuresLen := len(failures)
	if failuresLen != 0 {
		failuresText := ""
		for idx, f := range failures {
			failuresText += f.Failure
			if idx < len(failures)-1 {
				failuresText += ", "
			}
		}
		t.Errorf("Found %d failures in the code: %s", failuresLen, failuresText)
	}
}

// AssertSuccess checks if given rule runs correctly with no failures.
func AssertSuccess(t *testing.T, code string, testingRule rule.Rule, args rule.Arguments) {
	AssertSuccessWithName(t, code, "testing.go", testingRule, args)
}

// AssertFailures checks if given rule runs correctly with failures.
func AssertFailures(t *testing.T, code string, testingRule rule.Rule, args rule.Arguments) {
	annotations := extractFailures(code)
	if annotations == nil {
		t.Errorf("No failure annotations found in the code")
	}

	var fileSet token.FileSet
	file, err := file.New("testing.go", []byte(stripAnnotations(code)), &fileSet)
	if err != nil {
		t.Errorf("Cannot parse testing file: %s", err.Error())
	}
	failures := testingRule.Apply(file, args)
	totalFailures := len(failures)
	if totalFailures == 0 {
		t.Errorf("No failures in the code")
	}

	expectedFailures := len(annotations)
	if totalFailures != expectedFailures {
		t.Errorf("Expecting %d failures but got %d", expectedFailures, totalFailures)
	}

	for idx, f := range failures {
		if f.Node != nil {
			f.Position = rule.ToFailurePosition(f.Node.Pos(), f.Node.End(), file)
			failures[idx] = f
		}
	}

	for key, val := range annotations {
		matched := false
		start := file.ToPosition(token.Pos(val.start))
		end := file.ToPosition(token.Pos(val.end))

		for _, f := range failures {
			fmt.Println("#####", f.Position.Start.String(), f.Position.End.String())
			if f.Position.Start.String() == start.String() && f.Position.End.String() == end.String() {
				matched = true
				break
			}
		}

		if !matched {
			t.Errorf(`Failure annotation "%s" did not match any of the rule failures`, key)
		}
	}
}
