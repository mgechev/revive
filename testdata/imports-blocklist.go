package fixtures

import (
	"github.com/full/match"  //  MATCH /should not use the following blocklisted import: "github.com/full/match"/
	"bithub.com/full/match"
	"github.com/full/matche"
	"wildcard/between" // MATCH /should not use the following blocklisted import: "wildcard/between"/
	"wildcard/pkg1/between" // MATCH /should not use the following blocklisted import: "wildcard/pkg1/between"/
	"wildcard/pkg1/pkg2/between" // MATCH /should not use the following blocklisted import: "wildcard/pkg1/pkg2/between"/
	"wildcard/backward" // MATCH /should not use the following blocklisted import: "wildcard/backward"/
	"wildcard/backward/pkg" // MATCH /should not use the following blocklisted import: "wildcard/backward/pkg"/
	"wildcard/backward/pkg/pkg1" // MATCH /should not use the following blocklisted import: "wildcard/backward/pkg/pkg1"/
	"wildcard/forward" // MATCH /should not use the following blocklisted import: "wildcard/forward"/
	"pkg/wildcard/forward" // MATCH /should not use the following blocklisted import: "pkg/wildcard/forward"/
	"pkg/pkg1/wildcard/forward" // MATCH /should not use the following blocklisted import: "pkg/pkg1/wildcard/forward"/
	"full" // MATCH /should not use the following blocklisted import: "full"/
	"github.com/partical/match/fully"
	"strings"
)
