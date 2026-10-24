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

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall an installed JDK or unlink a linked JDK",
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

		var targetPath string
		var targetName string
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if name == strings.ToLower(target) || (strings.HasPrefix(name, vendor+"-") && strings.Contains(name, ver)) {
				targetPath = filepath.Join(cfg.VersionsDir, e.Name())
				targetName = e.Name()
				break
			}
		}

		if targetPath == "" {
			color.Red("Error: Version '%s' is not installed or linked.", target)
			return nil
		}

		// Check if active symlink points to this
		if link, err := os.Readlink(cfg.CurrentLink); err == nil && (link == targetPath || filepath.Base(link) == targetName) {
			_ = os.Remove(cfg.CurrentLink)
			color.Yellow("Note: Removing active default symlink.")
		}

		fileInfo, err := os.Lstat(targetPath)
		if err != nil {
			return err
		}

		// If it's a symlink (i.e. created with jvm link), remove ONLY the symlink
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("failed to unlink: %w", err)
			}
			color.Green("✓ Successfully unlinked '%s' (original external JDK was untouched)", targetName)
			return nil
		}

		// Otherwise, delete the full directory
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to remove version directory: %w", err)
		}

		color.Green("✓ Successfully uninstalled %s", targetName)
		return nil
	},
}
