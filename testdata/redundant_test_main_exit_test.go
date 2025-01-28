package fixtures

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	i := m.Run()
	teardown()
	// must not match because the go version of this module is less than 1.15
	os.Exit(i)
}

func setup() {
	fmt.Println("Setup")
}

func teardown() {
	fmt.Println("Teardown")
}

func Test_function(t *testing.T) {
	t.Error("Fail")
}
