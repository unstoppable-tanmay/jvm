package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
	"github.com/unstoppable-tanmay/jvm/internal/version"
)

var linkCmd = &cobra.Command{
	Use:   "link <path> [name]",
	Short: "Link an existing JDK directory into jvm (e.g. from Homebrew or /Library/Java)",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]
		absPath, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", absPath)
		}

		color.Cyan("🔍 Inspecting JDK at: %s...", absPath)
		info, err := version.InspectJDK(absPath)
		if err != nil {
			return fmt.Errorf("valid Java installation not found in %s: %w", absPath, err)
		}

		name := info.SuggestedID
		if len(args) == 2 {
			name = args[1]
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if err := cfg.EnsureDirectories(); err != nil {
			return err
		}

		destLink := filepath.Join(cfg.VersionsDir, name)
		if _, err := os.Lstat(destLink); err == nil {
			color.Yellow("⚠️  Version/link '%s' already exists in %s", name, cfg.VersionsDir)
			return nil
		}

		// Create symbolic link to the actual JAVA_HOME
		if err := os.Symlink(info.JavaHome, destLink); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}

		color.Green("✓ Successfully linked '%s' -> %s", name, info.JavaHome)
		color.HiWhite("\nTo use this version:")
		color.Yellow("  jvm use %s\n", name)

		return nil
	},
}
