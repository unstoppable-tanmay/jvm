package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unstoppable-tanmay/jvm/internal/config"
)

func TestConfigDefaultsAndCustomDir(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "jvm_home_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	// Test custom JVM_DIR override
	customDir := filepath.Join(tempHome, "custom_jvm")
	t.Setenv("JVM_DIR", customDir)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.RootDir != customDir {
		t.Errorf("expected RootDir %s, got %s", customDir, cfg.RootDir)
	}

	if err := cfg.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories failed: %v", err)
	}

	if _, err := os.Stat(cfg.VersionsDir); os.IsNotExist(err) {
		t.Errorf("VersionsDir was not created: %s", cfg.VersionsDir)
	}

	// Test saving settings (custom storage dir)
	newStorage := filepath.Join(tempHome, "external_storage")
	cfg.Settings.StorageDir = newStorage
	if err := cfg.SaveSettings(); err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	// Reload and verify
	reloaded, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("reload config failed: %v", err)
	}
	if reloaded.VersionsDir != newStorage {
		t.Errorf("expected reloaded VersionsDir %s, got %s", newStorage, reloaded.VersionsDir)
	}
}
