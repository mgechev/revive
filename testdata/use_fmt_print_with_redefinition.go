package fixtures

import (
	"fmt"
)

func print()   {}
func println() {}

type useFmtPrintT struct{}

func (useFmtPrintT) print(s string)   {}
func (useFmtPrintT) println(s string) {}

func useFmtPrint() {
	fmt.Println("just testing")
	fmt.Print("just testing")
	t := useFmtPrintT{}
	t.print("just testing")
	t.println("just testing")

	println("just testing")
	print("just testing")
}
