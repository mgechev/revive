package fixtures

func foo() {
	for i, newg := range groups {
		go func() { // MATCH /range value 'newg' seems to be referenced inside the closure/
			<-m.block
			newg.run(m.opts.Context)
		}()
		go func() {
			<-m.block
			newg.run(m.opts.Context)
		}(newg)
		go func() { // MATCH /range value 'i' seems to be referenced inside the closure/
			i++
			<-m.block
			newg.run(m.opts.Context)
		}(newg)
	}
}

func bar() {
	// from github.com/bazelbuild/bazel-gazelle/cmd/gazelle/update-repos.go:148:6
	for i, imp := range c.importPaths {
		go func(i int) { // MATCH /range value 'imp' seems to be referenced inside the closure/
			defer wg.Done()
			repo, err := repos.UpdateRepo(rc, imp)
			if err != nil {
				errs[i] = err
				return
			}
			repo.Remote = "" // don't set these explicitly
			repo.VCS = ""
			rule := repos.GenerateRule(repo)
			genRules[i] = rule
		}(i)
	}
}
