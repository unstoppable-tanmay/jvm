package platform

import (
	"fmt"
	"runtime"
)

// Info holds normalized operating system and architecture for JDK matching
type Info struct {
	OS          string // "macos", "linux", "windows"
	Arch        string // "aarch64", "x64"
	ArchiveType string // "tar.gz", "zip"
}

// Detect returns normalized platform info matching Foojay Disco API conventions
func Detect() (Info, error) {
	var osName, archName, archiveType string

	switch runtime.GOOS {
	case "darwin":
		osName = "macos"
		archiveType = "tar.gz"
	case "linux":
		osName = "linux"
		archiveType = "tar.gz"
	case "windows":
		osName = "windows"
		archiveType = "zip"
	default:
		return Info{}, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "arm64":
		archName = "aarch64"
	case "amd64":
		archName = "x64"
	default:
		return Info{}, fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	return Info{
		OS:          osName,
		Arch:        archName,
		ArchiveType: archiveType,
	}, nil
}
