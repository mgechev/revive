package cli

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	os.Unsetenv("HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	AppFs = afero.NewMemMapFs()
	os.Exit(m.Run())
}

func TestXDGConfigDirIsPrefferedFirst(t *testing.T) {
	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	xdgDirPath := "/tmp-iofs/xdg/config"
	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(xdgDirPath, 0755)
	AppFs.MkdirAll(homeDirPath, 0755)

	afero.WriteFile(AppFs, xdgDirPath+"/revive.toml", []byte("\n"), 0644)
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	afero.WriteFile(AppFs, homeDirPath+"/revive.toml", []byte("\n"), 0644)
	t.Setenv("HOME", homeDirPath)

	got := buildDefaultConfigPath()
	want := xdgDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestHomeConfigDir(t *testing.T) {

	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)

	afero.WriteFile(AppFs, homeDirPath+"/revive.toml", []byte("\n"), 0644)
	t.Setenv("HOME", homeDirPath)

	got := buildDefaultConfigPath()
	want := homeDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDir(t *testing.T) {
	xdgDirPath := "/tmp-iofs/xdg/config"
	AppFs.MkdirAll(xdgDirPath, 0755)

	afero.WriteFile(AppFs, xdgDirPath+"/revive.toml", []byte("\n"), 0644)
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := xdgDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDirNoFile(t *testing.T) {
	xdgDirPath := "/tmp-iofs/xdg/config"
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := xdgDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
