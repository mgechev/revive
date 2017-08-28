package main

import (
	"fmt"

	"github.com/mgechev/revive/defaultrules"
	"github.com/mgechev/revive/formatters"
	"github.com/mgechev/revive/linter"
	"github.com/mgechev/revive/rules"
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
