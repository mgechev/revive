package fixtures

func bare1() (int, int, error) {
	go func() (a int) {
		return // MATCH /avoid using bare returns, please add return expressions/
	}(5)
}

func bare2(a, b int) (int, error, int) {
	defer func() (a int) {
		return // MATCH /avoid using bare returns, please add return expressions/
	}(5)
}

func bare3(a string, b int) (a int, b float32, c string, d string) {
	go func() (a int, b int) {
		return a, b
	}(5, 6)

	defer func() (a int) {
		return a
	}(5)

	return // MATCH /avoid using bare returns, please add return expressions/
}

func bare4(a string, b int) string {
	return a
}

func bare5(a string, b int) {
	return
}
