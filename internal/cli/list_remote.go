package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/disco"
	"github.com/unstoppable-tanmay/jvm/internal/platform"
)

var (
	remoteVendor string
	remoteVer    string
)

var listRemoteCmd = &cobra.Command{
	Use:     "list-remote",
	Aliases: []string{"ls-remote"},
	Short:   "List available remote JDK versions from upstream repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		plat, err := platform.Detect()
		if err != nil {
			return err
		}

		client := disco.NewClient()
		color.Cyan("\nFetching available JDK releases for %s/%s...\n\n", plat.OS, plat.Arch)

		packages, err := client.QueryPackages(remoteVer, remoteVendor, plat)
		if err != nil {
			return err
		}

		if len(packages) == 0 {
			color.Yellow("No packages found matching the criteria.")
			return nil
		}

		color.HiWhite("%-15s %-15s %-10s %-25s", "DISTRIBUTION", "VERSION", "SUPPORT", "FILENAME")
		color.HiBlack("--------------------------------------------------------------------------------")

		// Display up to 25 packages
		maxShow := 25
		if len(packages) < maxShow {
			maxShow = len(packages)
		}

		for i := 0; i < maxShow; i++ {
			p := packages[i]
			support := p.TermOfSupport
			if support == "" {
				support = "standard"
			}
			fmt.Printf("%-15s %-15s %-10s %-25s\n", p.Distribution, p.JavaVersion, support, p.Filename)
		}

		if len(packages) > maxShow {
			color.HiBlack("\n... and %d more releases. Filter with --vendor or --version.", len(packages)-maxShow)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	listRemoteCmd.Flags().StringVarP(&remoteVendor, "vendor", "v", "temurin", "JDK Distribution/Vendor (e.g. temurin, corretto, zulu)")
	listRemoteCmd.Flags().StringVarP(&remoteVer, "version", "r", "", "Major version filter (e.g. 17, 21)")
}
