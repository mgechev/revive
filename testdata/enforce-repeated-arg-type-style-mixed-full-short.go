package fixtures

func compliantFunc(a int, b int, c string) (x, y int, z string) // Must not match - compliant with rule

func nonCompliantFunc1(a int, b int, c string) (x int, y int, z string) { panic("implement me") } // MATCH /repeated return type can be omitted/
func nonCompliantFunc2(a, b int, c string) (x, y int, z string)         { panic("implement me") } // MATCH /argument types should not be omitted/
