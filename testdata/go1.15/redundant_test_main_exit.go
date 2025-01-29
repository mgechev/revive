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
	// must not match because this is not a test file
	os.Exit(i)
}

func setup() {
	fmt.Println("Setup")
}

func teardown() {
	fmt.Println("Teardown")
}
