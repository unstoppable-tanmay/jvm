package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all locally installed and linked JDK versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		if _, err := os.Stat(cfg.VersionsDir); os.IsNotExist(err) {
			color.Yellow("No JDK versions installed yet. Try running: jvm install 21")
			return nil
		}

		entries, err := os.ReadDir(cfg.VersionsDir)
		if err != nil {
			return fmt.Errorf("failed to read versions directory: %w", err)
		}

		var currentTarget string
		if link, err := os.Readlink(cfg.CurrentLink); err == nil {
			currentTarget = filepath.Base(link)
		}

		installedCount := 0
		color.HiWhite("\nInstalled & Linked Java Versions:")
		color.HiBlack("--------------------------------")

		for _, e := range entries {
			fullPath := filepath.Join(cfg.VersionsDir, e.Name())
			fileInfo, err := os.Lstat(fullPath)
			if err != nil {
				continue
			}

			// Must be directory or symlink
			isSymlink := fileInfo.Mode()&os.ModeSymlink != 0
			if !fileInfo.IsDir() && !isSymlink {
				continue
			}

			installedCount++
			isCurrent := (e.Name() == currentTarget)

			var linkTarget string
			if isSymlink {
				if target, err := os.Readlink(fullPath); err == nil {
					linkTarget = target
				}
			}

			prefix := "   "
			if isCurrent {
				prefix = "-> "
			}

			extra := ""
			if isCurrent && linkTarget != "" {
				extra = fmt.Sprintf("(current default, linked: %s)", linkTarget)
			} else if isCurrent {
				extra = "(current default)"
			} else if linkTarget != "" {
				extra = fmt.Sprintf("(linked: %s)", linkTarget)
			}

			if isCurrent {
				color.Green("%s%-25s %s", prefix, e.Name(), extra)
			} else if linkTarget != "" {
				fmt.Printf("%s%-25s ", prefix, e.Name())
				color.HiBlack("%s\n", extra)
			} else {
				fmt.Printf("%s%-25s %s\n", prefix, e.Name(), extra)
			}
		}

		if installedCount == 0 {
			color.Yellow("No JDK versions installed yet. Try running: jvm install 21")
		} else {
			fmt.Println()
		}

		return nil
	},
}
