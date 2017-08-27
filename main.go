package main

import (
	"fmt"

	"github.com/mgechev/golinter/defaultrules"
	"github.com/mgechev/golinter/formatters"
	"github.com/mgechev/golinter/linter"
	"github.com/mgechev/golinter/rules"
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
  `

	linter := linter.New(func(file string) ([]byte, error) {
		return []byte(src), nil
	})
	var result []rules.Rule
	result = append(result, &defaultrules.LintElseRule{})

	failures, err := linter.Lint([]string{"foo.go", "bar.go", "baz.go"}, result)
	if err != nil {
		panic(err)
	}

	var formatter formatters.CLIFormatter
	output, err := formatter.Format(failures)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
