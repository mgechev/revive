package fixtures

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	i := m.Run()
	teardown()
	os.Exit(i)      // MATCH /redundant call to os.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
	syscall.Exit(i) // MATCH /redundant call to syscall.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
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
