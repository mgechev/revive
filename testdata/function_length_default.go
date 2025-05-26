package fixtures

func funLengthA() (a int) { // MATCH /maximum number of statements per function exceeded; max 50 but got 51/
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
	println()
}

func funLengthB(file *ast.File, fset *token.FileSet, lineLimit, stmtLimit int) []Message { // MATCH /maximum number of lines per function exceeded; max 75 but got 76/
	if true {
		a = b
		if false {
			c = d
			for _, f := range list {
				_, ok := f.(int64)
				if !ok {
					continue
				}
			}
		}
	}
	if true {
		a = b
		if false {
			c = d
			for _, f := range list {
				_, ok := f.(int64)
				if !ok {
					continue
				}
			}
			switch a {
			case 1:
				println()
			case 2:
				println()
				println()
			default:
				println()

			}
		}
	}
	if true {
		a = b
		if false {
			c = d
			for _, f := range list {
				_, ok := f.(int64)
				if !ok {
					continue
				}
			}
			switch a {
			case 1:
				println()
			case 2:
				println()
				println()
			default:
				println()

			}
		}
	}
	if true {
		a = b
		if false {
			c = d
			for _, f := range list {
				_, ok := f.(int64)
				if !ok {
					continue
				}
			}
			switch a {
			case 1:
				println()
			default:
				println()

			}
		}
	}
	return
}
