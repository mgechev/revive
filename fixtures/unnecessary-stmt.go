package fixtures

func foo(a, b, c, d int) {
	switch n := node.(type) { // MATCH /switch with only one case can be replaced by an if-then/
	case *ast.SwitchStmt:
		caseSelector := func(n ast.Node) bool {
			_, ok := n.(*ast.CaseClause)
			return ok
		}
		cases := pick(n.Body, caseSelector, nil)
		if len(cases) == 1 {
			cs, ok := cases[0].(*ast.CaseClause)
			if ok && len(cs.List) == 1 {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       n,
					Category:   "style",
					Failure:    "switch can be replaced by an if-then",
				})
			}
		}
	}
}

func bar() {
	a := 1

	switch a {
	case 1, 2:
		a++
	}

loop:
	for {
		switch a {
		case 1:
			a++
			println("one")
			break // MATCH /omit unnecessary break at the end of case clause/
		case 2:
			println("two")
			break loop
		default:
			println("default")
		}
	}

	return // MATCH /omit unnecessary return statement/
}
