package fixtures

import (
	"log"
	"os"
	"time"
)

func ur1() {
	ur1() // MATCH /unconditional recursive call/
}

func ur1bis() {
	if true {
		print()
	} else {
		switch {
		case true:
			println()
		default:
			for i := 0; i < 10; i++ {
				print()
			}
		}

	}

	ur1bis() // MATCH /unconditional recursive call/
}

func ur2tris() {
	for {
		println()
		ur2tris() // MATCH /unconditional recursive call/
	}
}

func ur2() {
	if true {
		return
	}

	ur2()
}

func ur3() {
	ur1()
}

func urn4() {
	if true {
		print()
	} else if false {
		return
	}

	ur4()
}

func urn5() {
	if true {
		return
	}

	if true {
		println()
	}

	ur5()
}

func ur2tris() {
	for true == false {
		println()
		ur2tris()
	}
}

type myType struct {
	foo int
	bar int
}

func (mt *myType) Foo() int {
	return mt.Foo() // MATCH /unconditional recursive call/
}

func (mt *myType) Bar() int {
	return mt.bar
}

func ur6() {
	switch {
	case true:
		return
	default:
		println()
	}

	ur6()
}

func ur7(a interface{}) {
	switch a.(type) {
	case int:
		return
	default:
		println()
	}

	ur7()
}

func ur8(a []int) {
	for _, i := range a {
		return
	}

	ur8()
}

func ur9(a []int) {
	for _, i := range a {
		ur9()
	}
}

func ur10() {
	select {
	case <-aChannel:
	case <-time.After(2 * time.Second):
		return
	}
	ur10()
}

func ur11() { // this pattern produces "infinite" number of goroutines
	go ur11() // MATCH /unconditional recursive call/
}

func ur12() {
	go foo(ur12())                   // MATCH /unconditional recursive call/
	go bar(1, "string", ur12(), 1.0) // MATCH /unconditional recursive call/
	go foo(bar())
}

func urn13() {
	if true {
		panic("")
	}
	urn13()
}

func urn14() {
	if true {
		os.Exit(1)
	}
	urn14()
}

func urn15() {
	if true {
		log.Panic("")
	}
	urn15()
}

func urn16(ch chan int) {
	for range ch {
		log.Panic("")
	}
	urn16(ch)
}

func urn17(ch chan int) {
	for range ch {
		print("")
	}
	urn17(ch) // MATCH /unconditional recursive call/
}

// Tests for #596
func (*fooType) BarFunc() {
	BarFunc()
}

func (_ *fooType) BazFunc() {
	BazFunc()
}

// Tests for #902
func falsePositiveFuncLiteral() {
	_ = foo(func() {
		falsePositiveFuncLiteral()
	})
}
func nr902() {
	go func() {
		nr902() // MATCH /unconditional recursive call/
	}()
}
