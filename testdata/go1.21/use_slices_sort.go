package fixtures

import "sort"

func useSlicesSort() {
	names := []string{}
	years := []int{}
	temperatures := []float64{}

	sort.Strings(names)                  // MATCH /replace sort.Strings by slices.Sort/
	sort.Ints(years)                     // MATCH /replace sort.Ints by slices.Sort/
	sort.Float64s(temperatures)          // MATCH /replace sort.Float64s by slices.Sort/
	sort.IntsAreSorted(years)            // MATCH /replace sort.IntsAreSorted by slices.IsSorted/
	sort.StringsAreSorted(names)         // MATCH /replace sort.StringsAreSorted by slices.IsSorted/
	sort.Float64sAreSorted(temperatures) // MATCH /replace sort.Float64sAreSorted by slices.IsSorted/

	sortable := sortable{}
	sort.Sort(sortable)                                             // MATCH /replace sort.Sort by slices.SortFunc/
	sort.Slice(years, func(i, j int) bool { return false })         // MATCH /replace sort.Slice by slices.SortFunc/
	sort.Stable(sortable)                                           // MATCH /replace sort.Stable by slices.SortStableFunc/
	sort.SliceStable(years, func(i, j int) bool { return false })   // MATCH /replace sort.SliceStable by slices.SortStableFunc/
	sort.IsSorted(sortable)                                         // MATCH /replace sort.IsSorted by slices.IsSortedFunc/
	sort.SliceIsSorted(years, func(i, j int) bool { return false }) // MATCH /replace sort.SliceIsSorted by slices.IsSortedFunc/
}

type sortable []int

func (sortable) Len() int           { return 0 }
func (sortable) Less(i, j int) bool { return true }
func (sortable) Swap(i, j int)      {}
