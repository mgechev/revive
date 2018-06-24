package fixtures

import (
	"syscall"
	"log"
	"os"
)

func foo0() {
    os.Exit(1) // MATCH /calls to os.Exit function should be made only in main() or init() functions/
}

func init() {
    log.Fatal("v ...interface{}")
}

func foo() {
    log.Fatalf(1) // MATCH /calls to log.Fatalf function should be made only in main() or init() functions/
}

func main() {
    log.Fatalln("v ...interface{}")
}

func bar() {
    log.Fatal(1) // MATCH /calls to log.Fatal function should be made only in main() or init() functions/
}

func bar2() {
    syscall.Exit(1) // MATCH /calls to syscall.Exit function should be made only in main() or init() functions/
}