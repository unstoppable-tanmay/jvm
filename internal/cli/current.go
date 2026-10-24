package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Display the active Java version and location",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		linkTarget, err := os.Readlink(cfg.CurrentLink)
		if err == nil {
			color.Green("Configured jvm default: %s", filepath.Base(linkTarget))
			color.HiBlack("Path: %s", linkTarget)
		} else {
			color.Yellow("No version set as jvm current default yet.")
		}

		// Also check the actual system java runtime
		javaBin, err := exec.LookPath("java")
		if err == nil {
			cmd := exec.Command(javaBin, "-version")
			output, _ := cmd.CombinedOutput()
			firstLine := strings.Split(string(output), "\n")[0]
			fmt.Printf("\nSystem active runtime: %s\n", strings.TrimSpace(firstLine))
			color.HiBlack("Binary: %s\n", javaBin)
		}

		return nil
	},
}
