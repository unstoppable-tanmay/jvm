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

var setDirCmd = &cobra.Command{
	Use:   "set-dir <path>",
	Short: "Set a custom JDK directory to use and scan for existing JDKs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]
		absPath, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			if err := os.MkdirAll(absPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		cfg.Settings.StorageDir = absPath
		if err := cfg.SaveSettings(); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		color.Green("✓ JDK storage directory set to: %s", absPath)

		// Automatically scan the provided path for any existing JDK installations and link them
		entries, err := os.ReadDir(absPath)
		if err == nil {
			linkedCount := 0
			for _, e := range entries {
				candidatePath := filepath.Join(absPath, e.Name())
				info, err := version.InspectJDK(candidatePath)
				if err != nil {
					continue
				}

				name := info.SuggestedID
				destLink := filepath.Join(cfg.RootDir, "versions", name)
				_ = os.MkdirAll(filepath.Dir(destLink), 0755)

				if _, err := os.Lstat(destLink); err == nil {
					continue
				}

				if err := os.Symlink(info.JavaHome, destLink); err == nil {
					color.Green("  ✓ Discovered and registered %s", name)
					linkedCount++
				}
			}
			if linkedCount > 0 {
				color.HiWhite("\nRegistered %d existing JDKs from %s", linkedCount, absPath)
			}
		}

		return nil
	},
}
