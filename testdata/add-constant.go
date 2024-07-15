package fixtures

import (
	"context"
	"fmt"
	"os"
)

type testLogger struct{}

func (l *testLogger) Info(ctx context.Context, msg string) {}

func getLogger() *testLogger {
	return &testLogger{}
}

func test1007() {
	getLogger().Info(context.Background(), "test1007")
	getLogger().Info(context.Background(), "test1007")
	getLogger().Info(context.Background(), "test1007")
}

func foo(a float32, b string, c any, d int) {
	a = 1.0 // ignore
	b = "ignore"
	c = 2              // ignore
	println("lit", 12) // MATCH /avoid magic numbers like '12', create a named constant for it/
	if a == 12.50 {    // MATCH /avoid magic numbers like '12.50', create a named constant for it/
		if b == "lit" {
			c = "lit" // MATCH /string literal "lit" appears, at least, 3 times, create a named constant for it/
		}
		for i := 0; i < 1; i++ {
			println("lit")
		}
	}

	println(0666)           // MATCH /avoid magic numbers like '0666', create a named constant for it/
	os.Chmod("test", 0666)  // ignore
	os.FindProcess(102100)  // ignore
	fmt.Println("test", 12) // ignore
	fmt.Printf("%d", 100)   // MATCH /avoid magic numbers like '100', create a named constant for it/
	myPrintln("%d", 100)    // MATCH /avoid magic numbers like '100', create a named constant for it/
	ignoredFunc(1000)       // ignore
	ignoredFunc1(1000)      // ignore - match regexp too

	println("The result of calling myFunc is: ", ignoredFunc(100))             // ignore
	println("result is: ", ignoredFunc(notIgnoredFunc(ignoredFunc(100))))      // ignore
	println("result of calling myFunc is: ", notIgnoredFunc(ignoredFunc(100))) // ignore

	println("result myFunc is: ", notIgnoredFunc(100))           // MATCH /avoid magic numbers like '100', create a named constant for it/
	println("The result is: ", ignoredFunc(notIgnoredFunc(100))) // MATCH /avoid magic numbers like '100', create a named constant for it/
}

func myPrintln(s string, num int) {

}

func ignoredFunc1(num int) int {
	return num
}

func ignoredFunc(num int) int {
	return num
}

func notIgnoredFunc(num int) int {
	return num
}

func tagsInStructLiteralsShouldBeOK() {
	a := struct {
		X int `json:"x"`
	}{}

	b := struct {
		X int `json:"x"`
	}{}

	c := struct {
		X int `json:"x"`
	}{}

	d := struct {
		X int `json:"x"`
		Y int `json:"y"`
	}{}

	e := struct {
		X int `json:"x"`
		Y int `json:"y"`
	}{}

	var f struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	var g struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}
