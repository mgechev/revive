package fixtures

func uselessFallthrough() {

	switch a {
	case 0:
		println()
		fallthrough
	default:
	}

	switch a {
	case 0:
		fallthrough // MATCH /this "fallthrough" can be removed by consolidating this case clause with the next one/
	case 1:
		println()
	default:
	}

	switch a {
	case 0:
		fallthrough // MATCH /this "fallthrough" can be removed by consolidating this case clause with the next one/
	case 1:
		fallthrough // MATCH /this "fallthrough" can be removed by consolidating this case clause with the next one/
	case 2:
		println()
	default:
	}

	switch a {
	case 0:
		fallthrough // json:{"MATCH": "this \"fallthrough\" can be removed by consolidating this case clause with the next one","Confidence": 0.8}
	default:
		println()
	}

	switch a {
	case 0:
		fallthrough // json:{"MATCH": "this \"fallthrough\" can be removed by consolidating this case clause with the next one","Confidence": 0.8}
	default:
		println()
	case 1:
		fallthrough // MATCH /this "fallthrough" can be removed by consolidating this case clause with the next one/
	case 2:
		println()
	}
}
