package pkg

var i interface{} // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/

type t interface{}   // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
type a = interface{} // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/

func any1(a interface{}) { // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
	m1 := map[interface{}]string{}     // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
	m2 := map[int]interface{}{}        // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
	a := []interface{}{}               // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
	m3 := make(map[int]interface{}, 1) // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
	a2 := make([]interface{}, 2)       // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/
}

func any2(a int) interface{} {} // MATCH /since GO 1.18 'interface{}' can be replaced by 'any'/

var ni interface{ Close() }

type nt interface{ Close() }
type na = interface{ Close() }

func nany1(a interface{ Close() }) {
	nm1 := map[interface{ Close() }]string{}
	nm2 := map[int]interface{ Close() }{}
	na := []interface{ Close() }{}
	nm3 := make(map[int]interface{ Close() }, 1)
	na2 := make([]interface{ Close() }, 2)
}

func nany2(a int) interface{ Close() } {}
