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
	fmt.Fprintf(nil, "") // MATCH /Unhandled error in call to function fmt.Fprintf/
	os.Chdir("..")
	_ = os.Chdir("..")
	return err
}
