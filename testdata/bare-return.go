package fixtures

func bare1() (int, int, error) {
	go func() (a int) {
		return // MATCH /avoid using bare returns, please add return expressions/
	}()
	return 1, 1, nil
}

func bare2(a, b int) (int, error, int) {
	defer func() (a int) {
		return // MATCH /avoid using bare returns, please add return expressions/
	}()
	return 1, nil, 1
}

func bare3(a string, b int) (ra int, rb float32, rc string, rd string) {
	go func() (a int, b int) {
		return a, b
	}()

	defer func() (a int) {
		return a
	}()

	return // MATCH /avoid using bare returns, please add return expressions/
}

func bare4(a string, b int) string {
	return a
}

func bare5(a string, b int) {
	return
}

// NR tests for issue #280
func f280_1() (err error) {
	func() {
		return
	}()

	return nil
}

func f280_2() (err error) {
	func() (r int) {
		return // MATCH /avoid using bare returns, please add return expressions/
	}()

	return nil
}

func f280_3() (err error) {
	func() (r int) {
		return 1
	}()

	return // MATCH /avoid using bare returns, please add return expressions/
}

func f280_4() (err error) {
	func() (r int) {
		return func() (r int) {
			return // MATCH /avoid using bare returns, please add return expressions/
		}()
	}()

	return nil
}
