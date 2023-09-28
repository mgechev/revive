package cli

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func init() {

    AppFs = afero.NewMemMapFs()
}

func TestHomeConfigDir(t *testing.T) {
	
    homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)

	afero.WriteFile(AppFs, homeDirPath+"/revive.toml", []byte("\n"), 0644)
	os.Setenv("HOME", homeDirPath)
	
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
	os.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := xdgDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestXDGConfigDirNoFile(t *testing.T) {

	xdgDirPath := "/tmp-iofs/xdg/config"
    os.Setenv("XDG_CONFIG_HOME", xdgDirPath)
    
    got := buildDefaultConfigPath()
	want := xdgDirPath + "/revive.toml"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
