package rule

import (
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// Deprecated: use ImportsBlocklistRule instead 
type ImportsBlacklistRule struct {
	blocklist []*regexp.Regexp
	sync.Mutex
}

func (r *ImportsBlacklistRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {

	var failures []lint.Failure

	return failures
}

func (*ImportsBlacklistRule) Name() string {
	return "imports-blacklist"
}
