package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	os.Unsetenv("HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	AppFs = afero.NewMemMapFs()
	homedir.DisableCache = true
	os.Exit(m.Run())
}

func TestXDGConfigDirIsPreferredFirst(t *testing.T) {
	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	homeDirPath := filepath.FromSlash("/tmp-iofs/home/tester")
	AppFs.MkdirAll(xdgDirPath, 0755)
	AppFs.MkdirAll(homeDirPath, 0755)

	afero.WriteFile(AppFs, filepath.Join(xdgDirPath, "revive.toml"), []byte("\n"), 0644)
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	afero.WriteFile(AppFs, filepath.Join(homeDirPath, "revive.toml"), []byte("\n"), 0644)
	t.Setenv("HOME", homeDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(xdgDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestHomeConfigDir(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	homeDirPath := filepath.FromSlash("/tmp-iofs/home/tester")
	AppFs.MkdirAll(homeDirPath, 0755)

	afero.WriteFile(AppFs, filepath.Join(homeDirPath, "revive.toml"), []byte("\n"), 0644)
	t.Setenv("HOME", homeDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(homeDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDir(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	AppFs.MkdirAll(xdgDirPath, 0755)

	afero.WriteFile(AppFs, filepath.Join(xdgDirPath, "revive.toml"), []byte("\n"), 0644)
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(xdgDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDirNoFile(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := ""

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		builtBy string
		want    string
	}{
		{
			name:    "Development version",
			version: defaultVersion,
			commit:  defaultCommit,
			date:    defaultDate,
			builtBy: defaultBuilder,
			want:    "version \n",
		},
		{
			name:    "Release version",
			version: "v1.5.0-12-g7ee4500-dev",
			commit:  "7ee4500e125e2d1b12653b2c8e140fec380919b4",
			date:    "2024-11-15 10:52 UTC",
			builtBy: "builder",
			want: `Version:	v1.5.0-12-g7ee4500-dev
Commit:		7ee4500e125e2d1b12653b2c8e140fec380919b4
Built		2024-11-15 10:52 UTC by builder
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVersion(tt.builtBy, tt.date, tt.commit, tt.version)

			if got != tt.want {
				t.Errorf("getVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
