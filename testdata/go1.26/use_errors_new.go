package pkg

import "fmt"

func errorsNew() (int, error) {
	fmt.Errorf("repo cannot be nil")
	errs := append(errs, fmt.Errorf("commit cannot be nil"))
	fmt.Errorf("unable to load base repo: %w", err)
	fmt.Errorf("Failed to get full commit id for origin/%s: %w", pr.BaseBranch, err)

	return 0, fmt.Errorf(msg + "something")
}
