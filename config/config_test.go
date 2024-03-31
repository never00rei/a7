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
	}{
		{"OnlyHOME", "/home/user", "", "/home/user/.config/.a7-journal"},
		{"OnlyXDG_CONFIG_HOME", "", "/home/user/.config", "/home/user/.config/.a7-journal"},
		{"BothSetPreferXDG_CONFIG_HOME", "/home/user", "/home/user/.config", "/home/user/.config/.a7-journal"},
		{"NeitherSet", "", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := config.BuildConfPath(tc.home, tc.xdgConfigHome)

			if got != tc.want {
				t.Errorf("config.BuildConfigPath() = %q, want %q", got, tc.want)
			}
		})
	}
}
