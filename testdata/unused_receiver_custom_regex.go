package fixtures

// all will ok with xxxParam if Arguments = [{allowRegex="^xxx"}]

func (xxxParam *SomeObj) f0() {}

// still works with _

func (_ *SomeObj) f1() {}

func (yyyParam *SomeObj) f2() { // MATCH /method receiver 'yyyParam' is not referenced in method's body, consider removing or renaming it to match ^xxx/
}
