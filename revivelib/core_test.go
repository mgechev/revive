package revivelib_test

import (
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
	"github.com/mgechev/revive/rule"
)

func TestReviveLint(t *testing.T) {
	// ARRANGE
	revive := getMockRevive(t)

	// ACT
	failures, err := revive.Lint(revivelib.Include("../testdata/if-return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ASSERT
	failureList := []lint.Failure{}

	for failure := range failures {
		failureList = append(failureList, failure)
	}

	const expected = 5

	got := len(failureList)
	if got != expected {
		t.Fatalf("Expected failures to have %d failures, but it has %d.", expected, got)
	}
}

func TestReviveFormat(t *testing.T) {
	// ARRANGE
	revive := getMockRevive(t)

	failuresChan, err := revive.Lint(revivelib.Include("../testdata/if-return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ACT
	color.NoColor = true
	failures, exitCode, err := revive.Format("stylish", failuresChan)
	// ASSERT
	if err != nil {
		t.Fatal(err)
	}

	errorMsgs := []string{
		"(91, 3)  https://revive.run/r#unreachable-code  unreachable code after this statement",
		"(98, 3)  https://revive.run/r#unreachable-code  unreachable code after this statement",
		"(15, 2)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
		"(88, 3)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
		"(95, 3)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
	}
	for _, errorMsg := range errorMsgs {
		if !strings.Contains(failures, errorMsg) {
			t.Fatalf("Expected formatted failures\n'%s'\nto contain\n'%s', but it didn't.", failures, errorMsg)
		}
	}

	const expected = 1
	if exitCode != expected {
		t.Fatalf("Expected exit code to be %d, but it was %d.", expected, exitCode)
	}
}

type mockRule struct{}

func (r *mockRule) Name() string {
	return "mock-rule"
}

func (r *mockRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	return nil
}

func getMockRevive(t *testing.T) *revivelib.Revive {
	t.Helper()

	conf, err := config.GetConfig("../defaults.toml")
	if err != nil {
		t.Fatal(err)
	}

	revive, err := revivelib.New(
		conf,
		true,
		2048,
		revivelib.NewExtraRule(&rule.IfReturnRule{}, lint.RuleConfig{}),
		revivelib.NewExtraRule(&mockRule{}, lint.RuleConfig{}),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	return revive
}
