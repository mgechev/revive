package fixtures

func compliantFunc(a, b int, c string) (x int, y int, z string) // Must not match - compliant with rule

func nonCompliantFunc1(a, b int, c string) (x, y int, z string)         { panic("implement me") } // MATCH /return types should not be omitted/
func nonCompliantFunc2(a int, b int, c string) (x int, y int, z string) { panic("implement me") } // MATCH /repeated argument type can be omitted/
