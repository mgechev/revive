package fixtures

import (
	"sync"
)

// Rule use-waitgroup-go shall not match because this file is a package with Go version < 1.25
func useWaitGroupGo() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		doSomething()
	}()

	wg.Add(1)
	go func() {
		doSomething()
		wg.Done()
	}()

	// from golang.org/x/tools/go/packages/packages.go/parseFiles
	for i, file := range filenames {
		wg.Add(1)
		go func(i int, filename string) {
			parsed[i], errors[i] = ld.parseFile(filename)
			wg.Done()
		}(i, file)
	}
	wg.Wait()

	// from kubernetes/pkg/kubelet/cm/devicemanager/manager_test.go/TestGetTopologyHintsWithUpdates
	// notice the rule spots a wg.Add(2) (vs wg.Add(1)) therefore using wg.Go is possible but requires
	// replacing the wg.Add and the next two go statements with two wg.Go
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < test.count; i++ {
			// simulate the device plugin to send device updates
			mimpl.genericDeviceUpdateCallback(testResourceName, devs)
		}
		updated.Store(true)
	}()
	go func() {
		defer wg.Done()
		for !updated.Load() {
			test.testfunc(mimpl)
		}
	}()
	wg.Wait()
}
