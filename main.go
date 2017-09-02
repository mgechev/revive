package main

import (
	"fmt"

	"github.com/mgechev/revive/defaultrule"
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
    // revive:disable no-else-return
    } else {
      return 23;
    }
    // revive:enable no-else-return
  }
  `

	linter := linter.New(func(file string) ([]byte, error) {
		return []byte(src), nil
	})
	var result []rule.Rule
	result = append(result, &defaultrule.LintElseRule{})

	failures, err := linter.Lint([]string{"foo.go" /*, "bar.go", "baz.go"*/}, result)
	if err != nil {
		panic(err)
	}

	var formatter formatter.CLIFormatter
	output, err := formatter.Format(failures)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
