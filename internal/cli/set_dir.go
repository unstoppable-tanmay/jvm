package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		cleanNewPath := filepath.Clean(absPath)

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		oldVersionsDir := filepath.Clean(cfg.VersionsDir)
		if cleanNewPath == oldVersionsDir {
			color.Yellow("Directory is already set to: %s", absPath)
			return nil
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			if err := os.MkdirAll(absPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}

		// Check if old versions directory has JDKs to migrate/symlink
		oldEntries, _ := os.ReadDir(oldVersionsDir)
		var existingJDKs []os.DirEntry
		for _, e := range oldEntries {
			full := filepath.Join(oldVersionsDir, e.Name())
			fi, err := os.Lstat(full)
			if err == nil && (fi.IsDir() || (fi.Mode()&os.ModeSymlink != 0)) {
				existingJDKs = append(existingJDKs, e)
			}
		}

		if len(existingJDKs) > 0 {
			color.Yellow("\nExisting JDKs found in %s (%d versions):", oldVersionsDir, len(existingJDKs))
			for _, e := range existingJDKs {
				fmt.Printf("  • %s\n", e.Name())
			}

			fmt.Println("\nHow would you like to handle existing JDKs?")
			fmt.Println("  [1] Move (migrate all existing JDKs to new directory)")
			fmt.Println("  [2] Link (symlink existing JDKs into new directory)")
			fmt.Println("  [3] Leave (keep them where they are without linking)")

			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("\nSelect an option [1-3]: ")
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)

				if choice == "1" {
					// Move
					moved := 0
					for _, e := range existingJDKs {
						src := filepath.Join(oldVersionsDir, e.Name())
						dst := filepath.Join(absPath, e.Name())
						if _, err := os.Lstat(dst); os.IsNotExist(err) {
							if err := os.Rename(src, dst); err == nil {
								moved++
							}
						}
					}
					color.Green("✓ Migrated %d JDKs to %s", moved, absPath)
					break
				} else if choice == "2" {
					// Symlink
					linked := 0
					for _, e := range existingJDKs {
						src := filepath.Join(oldVersionsDir, e.Name())
						dst := filepath.Join(absPath, e.Name())
						// Resolve target if it's already a symlink
						realTarget := src
						if linkTarget, err := os.Readlink(src); err == nil {
							realTarget = linkTarget
						}
						if _, err := os.Lstat(dst); os.IsNotExist(err) {
							if err := os.Symlink(realTarget, dst); err == nil {
								linked++
							}
						}
					}
					color.Green("✓ Linked %d existing JDKs into %s", linked, absPath)
					break
				} else if choice == "3" {
					color.HiBlack("Left existing JDKs in %s", oldVersionsDir)
					break
				}
				color.Red("Invalid selection. Please choose 1, 2, or 3.")
			}
		}

		cfg.Settings.StorageDir = absPath
		if err := cfg.SaveSettings(); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		color.Green("\n✓ JDK storage directory set to: %s", absPath)

		// Also scan the target directory for pre-existing non-jvm JDKs
		entries, err := os.ReadDir(absPath)
		if err == nil {
			discoveredCount := 0
			for _, e := range entries {
				candidatePath := filepath.Join(absPath, e.Name())
				if _, err := version.InspectJDK(candidatePath); err != nil {
					continue
				}
				discoveredCount++
			}
			if discoveredCount > 0 {
				color.HiBlack("Total recognized JDKs in %s: %d", absPath, discoveredCount)
			}
		}

		return nil
	},
}
