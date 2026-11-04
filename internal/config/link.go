package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// SetCurrentVersion updates ~/.jvm/current to point to targetPath.
// On Windows, if standard os.Symlink fails (e.g. Developer Mode disabled or no Admin privileges),
// it falls back to creating a Directory Junction via `cmd /c mklink /J <link> <target>`
// or writing a target path file fallback.
func SetCurrentVersion(cfg *Config, targetPath string) error {
	_ = os.Remove(cfg.CurrentLink)
	// On Windows, junctions might require removing directory via os.Remove
	_ = os.RemoveAll(cfg.CurrentLink)

	err := os.Symlink(targetPath, cfg.CurrentLink)
	if err == nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		// Try directory junction: mklink /J <link> <target>
		cmd := exec.Command("cmd", "/c", "mklink", "/J", cfg.CurrentLink, targetPath)
		if out, junctionErr := cmd.CombinedOutput(); junctionErr == nil {
			return nil
		} else {
			_ = out
		}

		// Fallback: write a text pointer ~/.jvm/current.path so tools can read it reliably
		targetFile := filepath.Join(cfg.RootDir, "current.path")
		if writeErr := os.WriteFile(targetFile, []byte(targetPath), 0644); writeErr == nil {
			return nil
		}
	}

	return err
}

// GetCurrentVersion reads the currently active target directory from cfg.CurrentLink.
func GetCurrentVersion(cfg *Config) (string, error) {
	// 1. Try reading symlink
	if link, err := os.Readlink(cfg.CurrentLink); err == nil {
		return link, nil
	}

	// 2. On Windows or fallback, check if current is a directory / junction
	if fi, err := os.Stat(cfg.CurrentLink); err == nil && fi.IsDir() {
		return cfg.CurrentLink, nil
	}

	// 3. Fallback: check current.path file
	currentPathFile := filepath.Join(cfg.RootDir, "current.path")
	if data, err := os.ReadFile(currentPathFile); err == nil {
		path := string(data)
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}
