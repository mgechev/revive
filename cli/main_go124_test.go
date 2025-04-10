//go:build go1.24

package cli

import "testing"

func TestGetDevelopmentVersion(t *testing.T) {
	got := getVersion(defaultBuilder, defaultDate, defaultCommit, defaultVersion)
	want := "version (devel)\n"

	if got != want {
		t.Errorf("getVersion() = %q, want %q", got, want)
	}
}
