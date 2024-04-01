package config_test

import (
	"testing"

	"github.com/never00rei/a7/config"
)

func TestBuildConfPath(t *testing.T) {
	tests := []struct {
		name          string
		home          string
		xdgConfigHome string
		want          string
		wantErr       bool
	}{
		{"OnlyHOME", "/home/user", "", "/home/user/.config/.a7-journal", false},
		{"OnlyXDG_CONFIG_HOME", "", "/home/user/.config", "/home/user/.config/.a7-journal", false},
		{"BothSetPreferXDG_CONFIG_HOME", "/home/user", "/home/user/.config", "/home/user/.config/.a7-journal", false},
		{"NeitherSet", "", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := config.BuildConfPath(tc.home, tc.xdgConfigHome)

			if tc.wantErr {
				if err == nil {
					t.Errorf("BuildConfPath() expected an error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("BuildConfPath() unexpected error = %v", err)
				}
			}

			if got != tc.want {
				t.Errorf("config.BuildConfigPath() = %q, want %q", got, tc.want)
			}
		})
	}
}
