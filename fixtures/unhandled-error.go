package fixtures

import (
	"fmt"
	"os"
)

func unhandledError1(a int) (int, error) {
	return a, nil
}

func unhandledError2() error {
	_, err := unhandledError1(1)
	unhandledError1(1)
	fmt.Fprintf(nil, "")
	os.Chdir("..")
	return err
}
