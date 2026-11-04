package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/unstoppable-tanmay/jvm/internal/disco"
	"github.com/unstoppable-tanmay/jvm/internal/platform"
)

var (
	remoteVendor string
	remoteVer    string
	remoteSearch string
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

		// If user passed a number in search (e.g. -s 8 or -s corretto,8) and didn't set -r,
		// use that number as remoteVer for the API query so the server filters to that major version
		// instead of downloading all 4600+ packages and doing naive string matching against build numbers.
		effectiveVer := remoteVer
		if effectiveVer == "" && remoteSearch != "" {
			terms := strings.FieldsFunc(remoteSearch, func(r rune) bool {
				return r == ',' || r == ' '
			})
			for _, t := range terms {
				clean := strings.TrimSpace(t)
				isNum := len(clean) > 0
				for _, ch := range clean {
					if ch < '0' || ch > '9' {
						isNum = false
						break
					}
				}
				if isNum {
					effectiveVer = clean
					break
				}
			}
		}

		packages, err := client.QueryPackages(effectiveVer, remoteVendor, plat)
		if err != nil {
			return err
		}

		// Filter packages by search terms if provided (supports comma-separated or space-separated terms)
		if remoteSearch != "" {
			terms := strings.FieldsFunc(remoteSearch, func(r rune) bool {
				return r == ',' || r == ' '
			})
			var filtered []disco.PackageItem
			for _, p := range packages {
				matchesAll := true
				targetText := strings.ToLower(fmt.Sprintf("%s %s %s %s", p.Distribution, p.JavaVersion, p.TermOfSupport, p.Filename))
				for _, term := range terms {
					cleanTerm := strings.ToLower(strings.TrimSpace(term))
					if cleanTerm == "" {
						continue
					}
					// If the search term is a pure integer, check major version match or java_version prefix
					isNum := true
					for _, ch := range cleanTerm {
						if ch < '0' || ch > '9' {
							isNum = false
							break
						}
					}
					if isNum {
						if fmt.Sprintf("%d", p.MajorVersion) != cleanTerm &&
							!strings.HasPrefix(p.JavaVersion, cleanTerm+".") &&
							!strings.HasPrefix(p.JavaVersion, cleanTerm+"+") &&
							p.JavaVersion != cleanTerm {
							matchesAll = false
							break
						}
					} else {
						if !strings.Contains(targetText, cleanTerm) {
							matchesAll = false
							break
						}
					}
				}
				if matchesAll {
					filtered = append(filtered, p)
				}
			}
			packages = filtered
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
			color.HiBlack("\n... and %d more releases. Filter with --vendor (-v), --version (-r), or --search (-s).", len(packages)-maxShow)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	listRemoteCmd.Flags().StringVarP(&remoteVendor, "vendor", "v", "", "JDK Distribution/Vendor filter (e.g. temurin, corretto, zulu, graalvm)")
	listRemoteCmd.Flags().StringVarP(&remoteVer, "version", "r", "", "Major version filter (e.g. 8, 17, 21)")
	listRemoteCmd.Flags().StringVarP(&remoteSearch, "search", "s", "", "Multi-keyword search filter (e.g. --search 'corretto,fx' or '8,lts')")
}
