package rule

import (
	"github.com/mgechev/revive/lint"
	"go/ast"
	"sync"
)

// A dumb set of potential failures
type failureCandidate struct {
	file *lint.File
	node ast.Node
}

// PkgLevelFailuresTestRule
type PkgLevelFailuresTestRule struct{}

var m sync.Mutex
var failureCandidates map[*lint.Package][]failureCandidate = map[*lint.Package][]failureCandidate{}

// Apply applies the rule to given file.
func (r *PkgLevelFailuresTestRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {

	// Here the code is a usual
	// for the sake of brevity, here we do not call ast.Walk to identify failures and/or collect data...
	//
	// We just collect arbitrary data
	m.Lock()
	failureCandidates[file.Pkg] = append(failureCandidates[file.Pkg], failureCandidate{file, file.AST.Name})
	m.Unlock()

	var failures []lint.Failure
	return failures
}

// Name returns the rule name.
func (r *PkgLevelFailuresTestRule) Name() string {
	return "package-level-failure-test"
}

// Reduce (implements Reducer interface)
func (r *PkgLevelFailuresTestRule) Reduce(p *lint.Package) []lint.PkgLevelFailure {
	// Here package level failures are calculated
	// For the POC we just pick the first candidate and create a failure with it
	if len(failureCandidates[p]) > 1 {
		f := failureCandidates[p][0]
		return []lint.PkgLevelFailure{lint.PkgLevelFailure{f.file, lint.Failure{
			Confidence: 1,
			Node:       f.node,
			Failure:    "this is a dumb package-level failure",
		}}}
	}

	return []lint.PkgLevelFailure{}

}
