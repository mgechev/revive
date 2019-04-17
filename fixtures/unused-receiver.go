package fixtures

import (
	"fmt"

	"github.com/mgechev/revive/lint"
)

func (f *Unix) Name() string { // MATCH /method receiver 'f' is not referenced in method's body, consider removing or renaming it as _/
	return "unix"
}

func (f *Unix) Format(failures <-chan lint.Failure, _ lint.RulesConfig) (string, error) { // MATCH /method receiver 'f' is not referenced in method's body, consider removing or renaming it as _/
	for failure := range failures {
		fmt.Printf("%v: [%s] %s\n", failure.Position.Start, failure.RuleName, failure.Failure)
	}
	return "", nil
}
