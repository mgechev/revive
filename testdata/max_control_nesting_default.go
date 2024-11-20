package fixtures

func mcn() {
	if true {
		if true {
			if true {
				if true {
					if true {
						if true { // MATCH /control flow nesting exceeds 5/
						}
					}
				}
			}
		}
	}
}
