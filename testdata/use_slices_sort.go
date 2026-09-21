package fixtures

import "sort"

func useSlicesSort() {
	names := []string{}
	years := []int{}
	temperatures := []float64{}

	sort.Strings(names)
	sort.Ints(years)
	sort.Float64s(temperatures)
	sort.IntsAreSorted(years)
	sort.StringsAreSorted(names)
	sort.Float64sAreSorted(temperatures)

	sortable := sortable{}
	sort.Sort(sortable)
	sort.Slice(years, func(i, j int) bool { return false })
	sort.Stable(sortable)
	sort.SliceStable(years, func(i, j int) bool { return false })
	sort.IsSorted(sortable)
	sort.SliceIsSorted(years, func(i, j int) bool { return false })
}

type sortable []int

func (sortable) Len() int           { return 0 }
func (sortable) Less(i, j int) bool { return true }
func (sortable) Swap(i, j int)      {}
