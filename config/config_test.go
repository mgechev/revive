package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestGetConfig(t *testing.T) {
	tt := map[string]struct {
		confPath   string
		wantConfig *lint.Config
		wantError  string
	}{
		"non-reg issue #470": {
			confPath:  "testdata/issue-470.toml",
			wantError: "",
		},
		"unknown file": {
			confPath:  "unknown",
			wantError: "cannot read the config file",
		},
		"malformed file": {
			confPath:  "testdata/malformed.toml",
			wantError: "cannot parse the config file",
		},
		"default config": {
			wantConfig: defaultConfig(),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := GetConfig(tc.confPath)
			switch {
			case err != nil && tc.wantError == "":
				t.Fatalf("Unexpected error\n\t%v", err)
			case err != nil && !strings.Contains(err.Error(), tc.wantError):
				t.Fatalf("Expected error\n\t%q\ngot:\n\t%v", tc.wantError, err)
			case tc.wantConfig != nil && reflect.DeepEqual(cfg, tc.wantConfig):
				t.Fatalf("Expected config\n\t%+v\ngot:\n\t%+v", tc.wantConfig, cfg)
			}

		})
	}
}
