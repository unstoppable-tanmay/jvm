package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a downloaded JDK version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		allCandidates, err := findCandidates(cfg.VersionsDir, target)
		if err != nil {
			return fmt.Errorf("failed to read versions directory: %w", err)
		}

		// Filter for actual installed directories (not external symlinks)
		var installedCandidates []candidateJDK
		var linkedCandidates []candidateJDK
		for _, c := range allCandidates {
			if c.IsSymlink {
				linkedCandidates = append(linkedCandidates, c)
			} else {
				installedCandidates = append(installedCandidates, c)
			}
		}

		if len(installedCandidates) == 0 {
			if len(linkedCandidates) > 0 {
				color.Yellow("Note: '%s' is an externally linked JDK.", target)
				color.Cyan("To unlink it without deleting original files, run:")
				color.HiWhite("  jvm unlink %s\n", target)
				return nil
			}
			color.Red("Error: Version '%s' is not installed.", target)
			return nil
		}

		selected := selectCandidate(installedCandidates, "uninstallation")
		if selected == nil {
			return nil
		}

		// Check if active symlink points to this
		if link, err := os.Readlink(cfg.CurrentLink); err == nil && (link == selected.Path || filepath.Base(link) == selected.Name) {
			_ = os.Remove(cfg.CurrentLink)
			color.Yellow("Note: Removing active default symlink.")
		}

		if err := os.RemoveAll(selected.Path); err != nil {
			return fmt.Errorf("failed to remove version directory: %w", err)
		}

		color.Green("✓ Successfully uninstalled %s", selected.Name)
		return nil
	},
}
