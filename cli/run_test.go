package cli

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// runReviveWith drives RunRevive with a clean flag set, since initConfig registers
// flags on the global CommandLine and would panic on a second registration.
func runReviveWith(t *testing.T, args ...string) int {
	t.Helper()

	oldArgs, oldCommandLine, oldUsage := os.Args, flag.CommandLine, flag.Usage
	t.Cleanup(func() {
		os.Args, flag.CommandLine, flag.Usage = oldArgs, oldCommandLine, oldUsage
		versionFlag, configPath, formatterName = false, "", ""
		setExitStatus, maxOpenFiles, excludePatterns = false, 0, nil
	})

	flag.CommandLine = flag.NewFlagSet("revive", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	excludePatterns = nil
	os.Args = append([]string{"revive"}, args...)

	return RunRevive()
}

// The point of returning the code instead of calling os.Exit: this test can exist at
// all. Before, a bad config reached fail() -> os.Exit(1), which took the test binary
// down with it, so no test could observe the failure path.
func TestRunReviveReturnsOneOnUnreadableConfig(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist.toml")

	if got := runReviveWith(t, "-config", missing, "./testdata/..."); got != 1 {
		t.Errorf("RunRevive() with an unreadable config = %d, want 1", got)
	}
}

func TestRunReviveReturnsZeroForVersion(t *testing.T) {
	if got := runReviveWith(t, "-version"); got != 0 {
		t.Errorf("RunRevive(-version) = %d, want 0", got)
	}
}

// The code must come from the lint result, not be a constant. With -set_exit_status,
// a file with findings exits 1 and a clean one exits 0 — so this pins the plumbing
// from revive.Format() through to the returned code.
func TestRunReviveReturnsLintExitStatus(t *testing.T) {
	dir := t.TempDir()

	write := func(name, src string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}

	clean := write("clean.go", "// Package clean is fine.\npackage clean\n\n// Add returns the sum of a and b.\nfunc Add(a, b int) int { return a + b }\n")
	// Missing package and exported-function comments: `exported` reports both.
	dirty := write("dirty.go", "package dirty\n\nfunc Add(a, b int) int { return a + b }\n")

	if got := runReviveWith(t, "-set_exit_status", clean); got != 0 {
		t.Errorf("RunRevive() on clean source = %d, want 0", got)
	}

	if got := runReviveWith(t, "-set_exit_status", dirty); got != 1 {
		t.Errorf("RunRevive() on source with findings = %d, want 1", got)
	}
}
