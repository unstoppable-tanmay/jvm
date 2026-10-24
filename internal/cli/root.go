package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is injected during build via ldflags (e.g. -ldflags "-X github.com/unstoppable-tanmay/jvm/internal/cli.Version=v1.0.0")
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "jvm",
	Short:   "jvm is a fast, lightweight Java Version Manager",
	Version: Version,
	Long: `jvm is a cross-platform Java Version Manager inspired by nvm.
It allows you to download, manage, link, and seamlessly switch between multiple JDK versions.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	envCmd.Hidden = true

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listRemoteCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(setDirCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(envCmd)
}
