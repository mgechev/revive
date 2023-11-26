package fixtures

import (
	"sync/atomic"
)

type Counter uint64

func AtomicTests() {
	x = atomic.AddUint64(&x, 1) // json:{"MATCH": "direct assignment to atomic value","Confidence": 1}
	x := uint64(1)

}
