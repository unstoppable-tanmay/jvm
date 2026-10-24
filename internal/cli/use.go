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

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific installed JDK version",
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
			return fmt.Errorf("failed to read installed versions: %w", err)
		}

		var matchedDir string
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if name == strings.ToLower(target) {
				matchedDir = filepath.Join(cfg.VersionsDir, e.Name())
				break
			}
			if strings.HasPrefix(name, vendor+"-") && strings.Contains(name, ver) {
				matchedDir = filepath.Join(cfg.VersionsDir, e.Name())
				break
			}
		}

		if matchedDir == "" {
			color.Red("Error: Version '%s' is not installed or linked.", target)
			color.Yellow("\nTo install it, run:")
			color.Cyan("  jvm install %s\n", target)
			return nil
		}

		// Update ~/.jvm/current symlink
		_ = os.Remove(cfg.CurrentLink)
		if err := os.Symlink(matchedDir, cfg.CurrentLink); err != nil {
			return fmt.Errorf("failed to link current version: %w", err)
		}

		versionName := filepath.Base(matchedDir)
		color.Green("✓ Now using %s", versionName)
		color.HiBlack("JAVA_HOME: %s\n", matchedDir)

		return nil
	},
}
