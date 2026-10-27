package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
	"github.com/unstoppable-tanmay/jvm/internal/version"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink <name>",
	Short: "Unlink an externally linked JDK without deleting original files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		vendor, ver := version.ParseIdentifier(target)

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		entries, err := os.ReadDir(cfg.VersionsDir)
		if err != nil {
			return fmt.Errorf("failed to read versions directory: %w", err)
		}

		var symlinkPath string
		var symlinkName string
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if name == strings.ToLower(target) || (strings.HasPrefix(name, vendor+"-") && strings.Contains(name, ver)) {
				fullPath := filepath.Join(cfg.VersionsDir, e.Name())
				fi, err := os.Lstat(fullPath)
				if err == nil && (fi.Mode()&os.ModeSymlink != 0) {
					symlinkPath = fullPath
					symlinkName = e.Name()
					break
				}
			}
		}

		if symlinkPath == "" {
			color.Red("Error: No linked JDK found matching '%s'.", target)
			color.Yellow("Use 'jvm uninstall' if you want to remove an installed JDK.")
			return nil
		}

		// Check if active symlink points to this
		if link, err := os.Readlink(cfg.CurrentLink); err == nil && (link == symlinkPath || filepath.Base(link) == symlinkName) {
			_ = os.Remove(cfg.CurrentLink)
			color.Yellow("Note: Removing active default symlink.")
		}

		if err := os.Remove(symlinkPath); err != nil {
			return fmt.Errorf("failed to unlink: %w", err)
		}

		color.Green("✓ Successfully unlinked '%s' (external JDK untouched)", symlinkName)
		return nil
	},
}
