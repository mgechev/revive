package fixtures

type gen1[T any] struct{}

func (g gen1[T]) f1() {}

func (g gen1[U]) f2() {}

func (n gen1[T]) f3() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

func (n gen1[U]) f4() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

func (n gen1[V]) f5() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

func (n *gen1[T]) f6() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

func (n *gen1[U]) f7() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

func (n *gen1[V]) f8() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen1/

type gen2[T1, T2 any] struct{}

func (g gen2[T1, T2]) f1() {}

func (g gen2[U1, U2]) f2() {}

func (n gen2[T1, T2]) f3() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/

func (n gen2[U1, U2]) f4() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/

func (n gen2[V1, V2]) f5() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/

func (n *gen2[T1, T2]) f6() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/

func (n *gen2[U1, U2]) f7() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/

func (n *gen2[V1, V2]) f8() {} // MATCH /receiver name n should be consistent with previous receiver name g for gen2/
