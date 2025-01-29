package fixtures

func compliantFunc(a, b int, c string) {}

func nonCompliantFunc1(a int, b int, c string) {}
func nonCompliantFunc2(a int, b, c int)        {}

type myStruct struct{}

func (m myStruct) compliantMethod(a, b int, c string) {}

func (m myStruct) nonCompliantMethod1(a int, b int, c string) {}
func (m myStruct) nonCompliantMethod2(a int, b, c int)        {}

func variadicFunction(a int, b ...int) {}

func singleArgFunction(a int) {}

func multiTypeArgs(a int, b string, c float64) {}

func mixedCompliance(a, b int, c int, d string) {}
