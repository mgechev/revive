package fixtures

// all will ok with xxxParam if Arguments = [{allowRegex="^xxx"}]

func f0(xxxParam int) {}

// still works with _

func f1(_ int) {}

func f2(yyyParam int) { // MATCH /parameter 'yyyParam' seems to be unused, consider removing or renaming it to match ^xxx/
}
