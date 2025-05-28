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

	println("just testing", something)   // MATCH /avoid using built-in function "println", replace it by "fmt.Fprintln(os.Stderr, "just testing", something)"/
	print("just testing", some, thing+1) // MATCH /avoid using built-in function "print", replace it by "fmt.Fprint(os.Stderr, "just testing", some, thing + 1)"/
}
