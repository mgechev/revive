package pkg

var i interface{}

type t interface{}
type a = interface{}

func any1(a interface{}) {
	m1 := map[interface{}]string{}
	m2 := map[int]interface{}{}
	a := []interface{}{}
	m3 := make(map[int]interface{}, 1)
	a2 := make([]interface{}, 2)
}

func any2(a int) interface{} {}

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
