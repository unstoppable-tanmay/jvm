package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
	"github.com/unstoppable-tanmay/jvm/internal/version"
)

var envShellType string

var envCmd = &cobra.Command{
	Use:    "env",
	Short:  "Output shell export commands for JAVA_HOME and PATH",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		// 1. Check for local .java-version or .jvmrc in current directory or parents
		currentWorkingDir, err := os.Getwd()
		var localVersion string
		if err == nil {
			checkDir := currentWorkingDir
			for {
				dotJavaVersion := filepath.Join(checkDir, ".java-version")
				if data, err := os.ReadFile(dotJavaVersion); err == nil {
					localVersion = strings.TrimSpace(string(data))
					break
				}
				dotJvmrc := filepath.Join(checkDir, ".jvmrc")
				if data, err := os.ReadFile(dotJvmrc); err == nil {
					localVersion = strings.TrimSpace(string(data))
					break
				}

				parent := filepath.Dir(checkDir)
				if parent == checkDir {
					break
				}
				checkDir = parent
			}
		}

		var targetDir string

		// If a local version is found, try to match it
		if localVersion != "" {
			vendor, ver := version.ParseIdentifier(localVersion)
			entries, err := os.ReadDir(cfg.VersionsDir)
			if err == nil {
				for _, e := range entries {
					name := strings.ToLower(e.Name())
					if name == strings.ToLower(localVersion) || (strings.HasPrefix(name, vendor+"-") && strings.Contains(name, ver)) {
						targetDir = filepath.Join(cfg.VersionsDir, e.Name())
						break
					}
				}
			}
		}

		// If not resolved from .java-version, fall back to ~/.jvm/current
		if targetDir == "" {
			if link, err := os.Readlink(cfg.CurrentLink); err == nil {
				targetDir = link
			}
		}

		if targetDir == "" {
			return nil
		}

		binDir := filepath.Join(targetDir, "bin")

		switch envShellType {
		case "powershell", "pwsh":
			fmt.Printf("$env:JAVA_HOME = %q\n$env:PATH = %q + [IO.Path]::PathSeparator + $env:PATH\n", targetDir, binDir)
		case "fish":
			fmt.Printf("set -gx JAVA_HOME %q;\nset -gx PATH %q $PATH;\n", targetDir, binDir)
		case "cmd":
			fmt.Printf("set \"JAVA_HOME=%s\" & set \"PATH=%s;%%PATH%%\"\n", targetDir, binDir)
		default:
			// POSIX (zsh, bash)
			fmt.Printf("export JAVA_HOME=%q;\nexport PATH=%q:\"$PATH\";\n", targetDir, binDir)
		}

		return nil
	},
}

func init() {
	envCmd.Flags().StringVar(&envShellType, "shell", "posix", "Target shell syntax (posix, powershell, fish)")
}
