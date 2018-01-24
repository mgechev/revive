package main

import (
	"fmt"
	"os"

	"github.com/mgechev/revive/formatter"
	"github.com/mgechev/revive/linter"
	"github.com/mgechev/revive/rule"
)

func main() {
	src := `
  package p

  func Test() {
    if true {
      return 42;
    } else {
      return 23;
    }
	}
	
	func foo_bar(a int, b int, c int, d int) {
		return a + b + c;
	}
  `

	revive := linter.New(func(file string) ([]byte, error) {
		return []byte(src), nil
	})
	var result []linter.Rule
	result = append(result, &rule.LintElseRule{}, &rule.ArgumentsLimitRule{}, &rule.NamesRule{})

	var config = linter.RulesConfig{
		"argument-limit": linter.RuleConfig{
			Arguments: []string{"3"},
			Severity:  linter.SeverityWarning,
		},
	}

	failures, err := revive.Lint([]string{"foo.go", "bar.go", "baz.go"}, result, config)
	if err != nil {
		panic(err)
	}

	formatChan := make(chan linter.Failure)
	exitChan := make(chan bool)
	var output string

	go (func() {
		var formatter formatter.CLIFormatter
		output, err = formatter.Format(formatChan, config)
		if err != nil {
			panic(err)
		}
		exitChan <- true
	})()

	exitCode := 0
	for f := range failures {
		if c, ok := config[f.RuleName]; ok && c.Severity == linter.SeverityError {
			exitCode = 1
		}
		formatChan <- f
	}
	close(formatChan)
	<-exitChan
	fmt.Println(output)

	os.Exit(exitCode)
}
