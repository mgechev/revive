package fixtures

func compliantFunc() (a, b int, c string) { panic("implement me") } // MATCH /return types should not be omitted/
func compliantFunc2() (int, int, string)  // Must not match - compliant with rule

func nonCompliantFunc1() (a int, b int, c string) { panic("implement me") } // Must not match - compliant with rule
func nonCompliantFunc2() (a int, b, c int)        { panic("implement me") } // MATCH /return types should not be omitted/

type myStruct struct{}

func (m myStruct) compliantMethod() (a, b int, c string) { panic("implement me") } // MATCH /return types should not be omitted/

func (m myStruct) nonCompliantMethod1() (a int, b int, c string) { panic("implement me") } // Must not match - compliant with rule
func (m myStruct) nonCompliantMethod2() (a int, b, c int)        { panic("implement me") } // MATCH /return types should not be omitted/

func singleArgFunction() (a int) { panic("implement me") } // Must not match - only one argument

func multiTypeArgs() (a int, b string, c float64) { panic("implement me") } // Must not match - different types for each argument

func mixedCompliance() (a, b int, c int, d string) { panic("implement me") } // MATCH /return types should not be omitted/
