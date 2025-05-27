package fixtures

import (
	"fmt"
)

type useFmtPrintT struct{}

func (useFmtPrintT) print(s string)   {}
func (useFmtPrintT) println(s string) {}

func useFmtPrint() {
	fmt.Println("just testing")
	fmt.Print("just testing")
	t := useFmtPrintT{}
	t.print("just testing")
	t.println("just testing")

	println("just testing") // MATCH /avoid using built-in function "println", use "fmt.Println" instead/
	print("just testing")   // MATCH /avoid using built-in function "print", use "fmt.Print" instead/
}
