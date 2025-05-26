package fixtures

import (
	"bithub.com/full/match"
	"full"                  // MATCH /should not use the following blocklisted import: "full"/
	"github.com/full/match" //  MATCH /should not use the following blocklisted import: "github.com/full/match"/
	"github.com/full/matche"
	"github.com/partical/match/fully"
	"pkg/pkg1/wildcard/forward" // MATCH /should not use the following blocklisted import: "pkg/pkg1/wildcard/forward"/
	"pkg/wildcard/forward"      // MATCH /should not use the following blocklisted import: "pkg/wildcard/forward"/
	"strings"
	"wildcard/backward"          // MATCH /should not use the following blocklisted import: "wildcard/backward"/
	"wildcard/backward/pkg"      // MATCH /should not use the following blocklisted import: "wildcard/backward/pkg"/
	"wildcard/backward/pkg/pkg1" // MATCH /should not use the following blocklisted import: "wildcard/backward/pkg/pkg1"/
	"wildcard/between"           // MATCH /should not use the following blocklisted import: "wildcard/between"/
	"wildcard/forward"           // MATCH /should not use the following blocklisted import: "wildcard/forward"/
	"wildcard/pkg1/between"      // MATCH /should not use the following blocklisted import: "wildcard/pkg1/between"/
	"wildcard/pkg1/pkg2/between" // MATCH /should not use the following blocklisted import: "wildcard/pkg1/pkg2/between"/
)
