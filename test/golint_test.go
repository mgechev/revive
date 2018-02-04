package test

import (
	"flag"
	"io/ioutil"
	"path"
	"regexp"
	"testing"

	"github.com/mgechev/revive/rule"

	"github.com/mgechev/revive/lint"
)

var lintMatch = flag.String("lint.match", "", "restrict fixtures matches to this pattern")

var rules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.NamesRule{},
	&rule.ElseRule{},
	&rule.IfReturnRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorsRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNameRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamesRule{},
	&rule.ContextKeyTypeRule{},
	&rule.ContextArgumentsRule{},
}

func TestAll(t *testing.T) {
	baseDir := "../fixtures/golint/"

	rx, err := regexp.Compile(*lintMatch)
	if err != nil {
		t.Fatalf("Bad -lint.match value %q: %v", *lintMatch, err)
	}

	fis, err := ioutil.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("ioutil.ReadDir: %v", err)
	}
	if len(fis) == 0 {
		t.Fatalf("no files in %v", baseDir)
	}
	for _, fi := range fis {
		if !rx.MatchString(fi.Name()) {
			continue
		}
		//t.Logf("Testing %s", fi.Name())
		src, err := ioutil.ReadFile(path.Join(baseDir, fi.Name()))
		if err != nil {
			t.Fatalf("Failed reading %s: %v", fi.Name(), err)
		}

		err = assertFailures(t, baseDir, fi, src, rules, map[string]lint.RuleConfig{})
		if err != nil {
			t.Errorf("Linting %s: %v", fi.Name(), err)
			continue
		}
	}
}
