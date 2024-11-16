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
