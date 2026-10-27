package cli

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/config"
	"github.com/unstoppable-tanmay/jvm/internal/disco"
	"github.com/unstoppable-tanmay/jvm/internal/downloader"
	"github.com/unstoppable-tanmay/jvm/internal/extractor"
	"github.com/unstoppable-tanmay/jvm/internal/platform"
	"github.com/unstoppable-tanmay/jvm/internal/version"
)

var (
	installURL    string
	installSHA256 string
)

var installCmd = &cobra.Command{
	Use:     "install [version]",
	Aliases: []string{"i"},
	Short:   "Download and install a specific JDK version (e.g. 21, temurin@17) or install from URL (--from)",
	Args:    cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if err := cfg.EnsureDirectories(); err != nil {
			return err
		}

		// Mode 1: Install from custom direct URL
		if installURL != "" {
			var customName string
			if len(args) > 0 {
				customName = args[0]
			}
			return installFromURL(cfg, installURL, installSHA256, customName)
		}

		// Mode 2: Standard install from Disco API
		if len(args) == 0 {
			return fmt.Errorf("please provide a version to install (e.g. 'jvm install 21') or use --from <url>")
		}

		rawInput := args[0]
		vendor, ver := version.ParseIdentifier(rawInput)

		plat, err := platform.Detect()
		if err != nil {
			return err
		}

		client := disco.NewClient()
		color.Cyan("🔍 Finding JDK for vendor=%s version=%s on %s/%s...", vendor, ver, plat.OS, plat.Arch)

		packages, err := client.QueryPackages(ver, vendor, plat)
		if err != nil {
			return err
		}

		if len(packages) == 0 {
			return fmt.Errorf("no JDK package found for %s (%s) on %s/%s", ver, vendor, plat.OS, plat.Arch)
		}

		pkg := packages[0]
		targetDirName := fmt.Sprintf("%s-%s", strings.ToLower(pkg.Distribution), pkg.JavaVersion)
		finalInstallPath := filepath.Join(cfg.VersionsDir, targetDirName)

		if _, err := os.Stat(finalInstallPath); err == nil {
			color.Yellow("⚠️  Version %s is already installed at: %s", targetDirName, finalInstallPath)
			return nil
		}

		color.Green("✓ Found package: %s (%s)", pkg.Filename, pkg.JavaVersion)
		color.Cyan("Fetching download details...")

		details, err := client.GetPackageDetails(pkg.ID)
		if err != nil {
			return err
		}

		downloadPath := filepath.Join(cfg.DownloadsDir, pkg.Filename)
		color.Cyan("⬇️  Downloading from upstream mirror...")
		if err := downloader.Download(details.DirectDownloadURI, downloadPath, details.Checksum); err != nil {
			return err
		}

		color.Cyan("📦 Extracting into %s...", finalInstallPath)
		tempExtractDir := filepath.Join(cfg.DownloadsDir, "temp_"+pkg.ID)
		_ = os.RemoveAll(tempExtractDir)
		if err := os.MkdirAll(tempExtractDir, 0755); err != nil {
			return err
		}
		defer os.RemoveAll(tempExtractDir)

		if err := extractor.Extract(downloadPath, tempExtractDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}

		realHome, err := version.LocateJavaHome(tempExtractDir)
		if err != nil {
			return err
		}

		if err := os.Rename(realHome, finalInstallPath); err != nil {
			return fmt.Errorf("failed to finalize installation directory: %w", err)
		}

		_ = os.Remove(downloadPath)

		color.Green("🎉 Successfully installed %s!", targetDirName)
		color.HiWhite("\nTo switch to this version now:")
		color.Yellow("  jvm use %s\n", targetDirName)

		return nil
	},
}

func installFromURL(cfg *config.Config, rawURL, expectedSHA, customName string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	filename := filepath.Base(u.Path)
	if filename == "" || filename == "/" || !strings.Contains(filename, ".") {
		filename = "downloaded_jdk.tar.gz"
	}

	downloadPath := filepath.Join(cfg.DownloadsDir, filename)
	color.Cyan("⬇️  Downloading JDK archive from URL:\n   %s", rawURL)

	if err := downloader.Download(rawURL, downloadPath, expectedSHA); err != nil {
		return err
	}

	tempExtractDir := filepath.Join(cfg.DownloadsDir, "temp_custom_url")
	_ = os.RemoveAll(tempExtractDir)
	if err := os.MkdirAll(tempExtractDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tempExtractDir)

	color.Cyan("📦 Extracting archive...")
	if err := extractor.Extract(downloadPath, tempExtractDir); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	info, err := version.InspectJDK(tempExtractDir)
	if err != nil {
		return fmt.Errorf("failed to detect Java runtime inside extracted archive: %w", err)
	}

	targetName := customName
	if targetName == "" {
		targetName = info.SuggestedID
	}

	finalInstallPath := filepath.Join(cfg.VersionsDir, targetName)
	if _, err := os.Stat(finalInstallPath); err == nil {
		color.Yellow("⚠️  Version/directory '%s' already exists in %s", targetName, cfg.VersionsDir)
		return nil
	}

	if err := os.Rename(info.JavaHome, finalInstallPath); err != nil {
		return fmt.Errorf("failed to finalize installation directory: %w", err)
	}

	_ = os.Remove(downloadPath)

	color.Green("🎉 Successfully installed custom JDK as '%s'!", targetName)
	color.HiWhite("\nTo switch to this version now:")
	color.Yellow("  jvm use %s\n", targetName)

	return nil
}

func init() {
	installCmd.Flags().StringVarP(&installURL, "from", "u", "", "Direct download URL for JDK archive (.tar.gz, .zip)")
	installCmd.Flags().StringVar(&installURL, "url", "", "Alias for --from")
	installCmd.Flags().StringVar(&installSHA256, "sha256", "", "Optional SHA256 checksum to verify downloaded archive")
}
