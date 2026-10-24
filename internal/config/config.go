package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UserSettings represents persistent preferences in ~/.jvm/config.json
type UserSettings struct {
	StorageDir string   `json:"storage_dir,omitempty"`
	ScanDirs   []string `json:"scan_dirs,omitempty"`
}

// Config holds core filesystem paths for jvm
type Config struct {
	RootDir      string
	VersionsDir  string
	DownloadsDir string
	CurrentLink  string
	ConfigFile   string
	Settings     UserSettings
}

// GetRootDir returns the base directory for jvm (default ~/.jvm or overridden by JVM_DIR)
func GetRootDir() (string, error) {
	if custom := os.Getenv("JVM_DIR"); custom != "" {
		return custom, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".jvm"), nil
}

// LoadConfig initializes and returns paths configuration, loading settings from config.json
func LoadConfig() (*Config, error) {
	rootDir, err := GetRootDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(rootDir, "config.json")
	var settings UserSettings

	if data, err := os.ReadFile(configFile); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	versionsDir := filepath.Join(rootDir, "versions")
	if settings.StorageDir != "" {
		versionsDir = settings.StorageDir
	}

	cfg := &Config{
		RootDir:      rootDir,
		VersionsDir:  versionsDir,
		DownloadsDir: filepath.Join(rootDir, "downloads"),
		CurrentLink:  filepath.Join(rootDir, "current"),
		ConfigFile:   configFile,
		Settings:     settings,
	}

	return cfg, nil
}

// SaveSettings writes UserSettings to config.json
func (c *Config) SaveSettings() error {
	data, err := json.MarshalIndent(c.Settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.ConfigFile, data, 0644)
}

// EnsureDirectories creates required directories if they don't exist
func (c *Config) EnsureDirectories() error {
	dirs := []string{c.RootDir, c.VersionsDir, c.DownloadsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
