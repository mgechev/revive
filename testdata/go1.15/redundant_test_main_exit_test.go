package fixtures

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse() // must not match
	if testing.Short() {
		m.Run()
		return
	}
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

func Test_os_exit(t *testing.T) {
	// must not match because this is not TestMain function
	os.Exit(1)
}

func Test_syscall_exit(t *testing.T) {
	// must not match because this is not TestMain function
	syscall.Exit(1)
}
