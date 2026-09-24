package a

import (
	"errors"
	"fmt"
)

var Public_var = 1 // want `don't use underscores in Go names; var Public_var should be PublicVar`

func doWork() error {
	fmt.Println("hello") // want `Unhandled error in call to function fmt\.Println`
	return errors.New(fmt.Sprintf("failed: %d", 42)) // want `should replace errors\.New\(fmt\.Sprintf\(\.\.\.\)\) with fmt\.Errorf\(\.\.\.\)`
}
