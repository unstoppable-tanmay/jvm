package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific installed JDK version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		candidates, err := findCandidates(cfg.VersionsDir, target)
		if err != nil {
			return fmt.Errorf("failed to read installed versions: %w", err)
		}

		if len(candidates) == 0 {
			color.Red("Error: Version '%s' is not installed or linked.", target)
			color.Yellow("\nTo install it, run:")
			color.Cyan("  jvm install %s\n", target)
			return nil
		}

		selected := selectCandidate(candidates, "switching active JDK")
		if selected == nil {
			return nil
		}

		// Update ~/.jvm/current symlink
		_ = os.Remove(cfg.CurrentLink)
		if err := os.Symlink(selected.Path, cfg.CurrentLink); err != nil {
			return fmt.Errorf("failed to link current version: %w", err)
		}

		versionName := filepath.Base(selected.Path)
		color.Green("✓ Now using %s", versionName)
		color.HiBlack("JAVA_HOME: %s\n", selected.Path)

		return nil
	},
}
