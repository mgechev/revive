package pkg

func datarace(i int) (r int) {
	var l int
	go func(ann int) {
		println(i, r, l, ann)
	}(1)
	i = 2
	l = 3
	r = 2
	return 1
}
