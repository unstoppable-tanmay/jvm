package platform_test

import (
	"runtime"
	"testing"

	"github.com/unstoppable-tanmay/jvm/internal/platform"
)

func TestDetect(t *testing.T) {
	info, err := platform.Detect()
	if err != nil {
		t.Fatalf("unexpected error on platform.Detect: %v", err)
	}

	switch runtime.GOOS {
	case "darwin":
		if info.OS != "macos" || info.ArchiveType != "tar.gz" {
			t.Errorf("expected macos and tar.gz, got %s and %s", info.OS, info.ArchiveType)
		}
	case "linux":
		if info.OS != "linux" || info.ArchiveType != "tar.gz" {
			t.Errorf("expected linux and tar.gz, got %s and %s", info.OS, info.ArchiveType)
		}
	case "windows":
		if info.OS != "windows" || info.ArchiveType != "zip" {
			t.Errorf("expected windows and zip, got %s and %s", info.OS, info.ArchiveType)
		}
	}

	switch runtime.GOARCH {
	case "arm64":
		if info.Arch != "aarch64" {
			t.Errorf("expected aarch64, got %s", info.Arch)
		}
	case "amd64":
		if info.Arch != "x64" {
			t.Errorf("expected x64, got %s", info.Arch)
		}
	}
}
