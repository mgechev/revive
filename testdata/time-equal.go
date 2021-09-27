package pkg

import "time"

func eqlOp() bool {
	t := time.Now()
	u := t
	return t == u // MATCH /use t.Equal(u) instead of "==" operator/
}

func eqlFun() bool {
	t := time.Now()
	u := t
	return t.Equal(t)
}

func neqlOp() bool {
	t := time.Now()
	u := t.Add(2 * time.Second)
	return t != u // MATCH /use !t.Equal(u) instead of "!=" operator/
}

func neqlFun() bool {
	t := time.Now()
	u := t.Sub(2 * time.Second)
	return !t.Equal(t)
}
