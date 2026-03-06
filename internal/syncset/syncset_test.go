package syncset_test

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/mgechev/revive/internal/syncset"
)

func TestNew_ConcurrentElementsAreInitiallyEmpty(t *testing.T) {
	t.Parallel()

	set := syncset.New()
	if set == nil {
		t.Fatal("New() returned nil")
	}

	const readers = 64
	var wg sync.WaitGroup

	errs := make(chan error, readers)
	for range readers {
		wg.Go(func() {
			if got := len(set.Elements()); got != 0 {
				errs <- fmt.Errorf("Elements() len = %d, want 0", got)
			}
		})
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

func TestSet_ZeroValueConcurrentUniqueAdds(t *testing.T) {
	t.Parallel()

	var set syncset.Set

	const workers = 128
	var wg sync.WaitGroup
	var failedAdds atomic.Int64

	for i := range workers {
		wg.Go(func() {
			if added := set.AddIfAbsent(fmt.Sprintf("value-%d", i)); !added {
				failedAdds.Add(1)
			}
		})
	}

	wg.Wait()

	if got := failedAdds.Load(); got != 0 {
		t.Errorf("AddIfAbsent() returned false %d times for unique values, want 0", got)
	}

	elements := set.Elements()
	if got := len(elements); got != workers {
		t.Errorf("Elements() len = %d, want %d", got, workers)
	}
}

func TestSet_AddIfAbsent_ConcurrentSameValueAddedOnce(t *testing.T) {
	t.Parallel()

	set := syncset.New()

	const workers = 256
	var wg sync.WaitGroup
	var additions atomic.Int64

	for range workers {
		wg.Go(func() {
			if set.AddIfAbsent("same-value") {
				additions.Add(1)
			}
		})
	}

	wg.Wait()

	if got := additions.Load(); got != 1 {
		t.Errorf("AddIfAbsent() reported %d additions, want 1", got)
	}

	want := []string{"same-value"}
	if got := set.Elements(); !reflect.DeepEqual(want, got) {
		t.Errorf("Elements() = %v, want %v", got, want)
	}
}
