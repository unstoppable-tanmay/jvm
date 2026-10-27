package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/unstoppable-tanmay/jvm/internal/version"
)

type candidateJDK struct {
	Name      string
	Path      string
	IsSymlink bool
	Target    string
}

func selectCandidate(candidates []candidateJDK, actionName string) *candidateJDK {
	if len(candidates) == 1 {
		return &candidates[0]
	}

	color.Yellow("\nMultiple matching JDKs found for %s:", actionName)
	for i, c := range candidates {
		tag := "installed"
		if c.IsSymlink {
			tag = fmt.Sprintf("linked -> %s", c.Target)
		}
		fmt.Printf("  [%d] %-30s (%s)\n", i+1, c.Name, tag)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nSelect an option [1-%d] or 'c' to cancel: ", len(candidates))
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil
		}
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "c" || input == "" {
			color.Yellow("Operation cancelled.")
			return nil
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= len(candidates) {
			return &candidates[choice-1]
		}
		color.Red("Invalid selection. Please enter a number between 1 and %d.", len(candidates))
	}
}

func findCandidates(versionsDir, target string) ([]candidateJDK, error) {
	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	vendor, ver := version.ParseIdentifier(target)
	var matches []candidateJDK

	for _, e := range entries {
		name := strings.ToLower(e.Name())
		fullPath := filepath.Join(versionsDir, e.Name())
		fi, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		isSymlink := (fi.Mode()&os.ModeSymlink != 0)
		if !fi.IsDir() && !isSymlink {
			continue
		}

		var targetDest string
		if isSymlink {
			if link, err := os.Readlink(fullPath); err == nil {
				targetDest = link
			}
		}

		if name == strings.ToLower(target) || (strings.HasPrefix(name, vendor+"-") && strings.Contains(name, ver)) {
			matches = append(matches, candidateJDK{
				Name:      e.Name(),
				Path:      fullPath,
				IsSymlink: isSymlink,
				Target:    targetDest,
			})
		}
	}

	return matches, nil
}
