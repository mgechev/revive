package fixtures

func datarace() (r int, c char) {
	for _, p := range []int{1, 2} {
		go func() {
			print(r) // MATCH /potential datarace: return value r is captured (by-reference) in goroutine/
			print(p) // MATCH /datarace: range value p is captured (by-reference) in goroutine/
		}()
		for i, p1 := range []int{1, 2} {
			a := p1
			go func() {
				print(r)  // MATCH /potential datarace: return value r is captured (by-reference) in goroutine/
				print(p)  // MATCH /datarace: range value p is captured (by-reference) in goroutine/
				print(p1) // MATCH /datarace: range value p1 is captured (by-reference) in goroutine/
				print(a)
				print(i) // MATCH /datarace: range value i is captured (by-reference) in goroutine/
			}()
			print(i)
			print(p)
			go func() {
				_ = c // MATCH /potential datarace: return value c is captured (by-reference) in goroutine/
			}()
		}
		print(p1)
	}
	go func() {
		print(r) // MATCH /potential datarace: return value r is captured (by-reference) in goroutine/
	}()
	print(r)
}
