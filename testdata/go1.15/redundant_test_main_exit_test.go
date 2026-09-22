package fixtures

import (
	"flag"
	"fmt"
	"log"
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

func TestMain(m *testing.M) {
	flag.Parse() // must not match
	if testing.Short() {
		m.Run()
		return
	}
}

var testBackend string

func TestMain(m *testing.M) {
	flag.StringVar(&testBackend, "backend", "", "selects backend")
	flag.Parse() // must not match
}

func TestMain(m *testing.M) {
	fs := flag.NewFlagSet("fs", flag.ExitOnError) // must not match
	fs.Parse()
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

func TestMain(m *testing.M) {
	os.Exit(m.Run()) // MATCH /redundant call to os.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
}

func TestMain(m *testing.M) {
	var code int
	code = m.Run()
	os.Exit((code)) // MATCH /redundant call to os.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
}

func TestMain(mm *testing.M) {
	code := mm.Run()
	syscall.Exit(code) // MATCH /redundant call to syscall.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
}

func TestMain(m *testing.M) {
	if err := setupErr(); err != nil {
		fmt.Println("setup failed:", err)
		os.Exit(1) // must not match because the exit status is not the result of m.Run
	}

	m.Run()

	if err := teardownErr(); err != nil {
		fmt.Println("teardown failed:", err)
		os.Exit(1) // must not match because the exit status is not the result of m.Run
	}
}

func TestMain(m *testing.M) {
	m.Run()
	syscall.Exit(3) // must not match because the exit status is not the result of m.Run
}

func TestMain(m *testing.M) {
	if err := setupErr(); err != nil {
		log.Fatal("setup failed") // must not match because log.Fatal always exits with status 1
	}
	m.Run()
}

func setupErr() error {
	return nil
}

func teardownErr() error {
	return nil
}

func TestMain(m *testing.M) {
	var code = m.Run()
	os.Exit(code) // MATCH /redundant call to os.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
}

func TestMain(m *testing.M) {
	os.Exit((m).Run()) // MATCH /redundant call to os.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15/
}

func TestMain(m *testing.M) {
	code := m.Run()
	if code == 0 && leaked() {
		code = 1
	}
	os.Exit(code) // must not match because the exit status is not the result of m.Run
}

func TestMain(m *testing.M) {
	code := m.Run()
	code++
	os.Exit(code) // must not match because the exit status is not the result of m.Run
}

func TestMain(m *testing.M) {
	code := m.Run()
	code |= 1
	syscall.Exit(code) // must not match because the exit status is not the result of m.Run
}

func TestMain(m *testing.M) {
	code := m.Run()
	adjust(&code)
	os.Exit(code) // must not match because the exit status is not the result of m.Run
}

func TestMain(m *testing.M) {
	if err := setupErr(); err != nil {
		code := 1
		os.Exit(code) // must not match because the exit status is not the result of m.Run
	}
	code := m.Run()
	os.Exit(code) // must not match because a variable with the same name is written elsewhere
}

func TestMain(m *testing.M) {
	code := 2
	if setupErr() != nil {
		os.Exit(code) // must not match because the exit status is not the result of m.Run
	}
	code = m.Run()
	os.Exit(code) // must not match because the variable is written more than once
}

type suite struct{}

func (suite) TestMain(m *testing.M) {
	os.Exit(m.Run()) // must not match because a method is never the test entry point
}

func leaked() bool {
	return false
}

func adjust(code *int) {
	*code = 1
}
