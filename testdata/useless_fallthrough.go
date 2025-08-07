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

	switch a {
	case 0:
		// This a comment on the case, the confidence must be 0.5
		fallthrough // json:{"MATCH": "this \"fallthrough\" can be removed by consolidating this case clause with the next one","Confidence": 0.5}
	default:
		println()
	}

	switch goos {
	case "linux":
		// TODO(bradfitz): be fancy and use linkat with AT_EMPTY_PATH to avoid
		// copying? I couldn't get it to work, though.
		// For now, just do the same thing as every other Unix and copy
		// the binary.
		fallthrough // json:{"MATCH": "this \"fallthrough\" can be removed by consolidating this case clause with the next one","Confidence": 0.5}
	case "darwin", "freebsd", "openbsd", "netbsd":
		return
	case "windows":
		return
	default:
		return
	}

	switch a {
	case 0:
		//foo:bar
		fallthrough // json:{"MATCH": "this \"fallthrough\" can be removed by consolidating this case clause with the next one","Confidence": 0.5}
	default:
		println()
	}

}
