package main

import (
	"fmt"

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
	
	func foobar(a int, b int, c int, d int) {
		return a + b + c;
	}
  `

	revive := linter.New(func(file string) ([]byte, error) {
		return []byte(src), nil
	})
	var result []linter.Rule
	result = append(result, &rule.LintElseRule{}, &rule.ArgumentsLimitRule{})

	var config = linter.RulesConfig{
		"argument-limit": []string{"3"},
	}

	failures, err := revive.Lint([]string{"foo.go", "bar.go", "baz.go"}, result, config)
	if err != nil {
		panic(err)
	}

	var formatter formatter.JSONFormatter
	output, err := formatter.Format(failures)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
