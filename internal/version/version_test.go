package version_test

import (
	"testing"

	"github.com/unstoppable-tanmay/jvm/internal/version"
)

func TestParseIdentifier(t *testing.T) {
	tests := []struct {
		input       string
		wantVendor  string
		wantVersion string
	}{
		{"21", "temurin", "21"},
		{"temurin@17", "temurin", "17"},
		{"corretto@21", "corretto", "21"},
		{"zulu-11", "zulu", "11"},
		{"graalvm@21", "graalvm", "21"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotVendor, gotVer := version.ParseIdentifier(tt.input)
			if gotVendor != tt.wantVendor || gotVer != tt.wantVersion {
				t.Errorf("ParseIdentifier(%q) = (%q, %q), want (%q, %q)",
					tt.input, gotVendor, gotVer, tt.wantVendor, tt.wantVersion)
			}
		})
	}
}
