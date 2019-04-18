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

func (u *Unix) Foo() (string, error) { // MATCH /method receiver 'u' is not referenced in method's body, consider removing or renaming it as _/
	for failure := range failures {
		u := 1 // shadowing the receiver
		fmt.Printf("%v\n", u)
	}
	return "", nil
}

func (u *Unix) Foos() (string, error) {
	for failure := range failures {
		u := 1 // shadowing the receiver
		fmt.Printf("%v\n", u)
	}

	return u, nil // use of the receiver
}

func (u *Unix) Bar() (string, error) {
	for failure := range failures {
		u.path = nil // modifies the receiver
		fmt.Printf("%v\n", u)
	}
	return "", nil
}
