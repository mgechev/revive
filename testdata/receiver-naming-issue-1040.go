package fixtures

func (o *o) f1()     {}
func (tw *tw) f2()   {}
func (thr *thr) f3() {} // MATCH /receiver name thr is longer than 2 characters/
