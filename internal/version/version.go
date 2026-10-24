package version

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// ParseIdentifier splits input into vendor and version.
// e.g.:
// "21" -> ("temurin", "21")
// "temurin@21" -> ("temurin", "21")
// "corretto-17" -> ("corretto", "17")
// "zulu@11" -> ("zulu", "11")
func ParseIdentifier(input string) (vendor string, ver string) {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "@") {
		parts := strings.SplitN(input, "@", 2)
		return strings.ToLower(parts[0]), parts[1]
	}
	if strings.Contains(input, "-") {
		parts := strings.SplitN(input, "-", 2)
		switch strings.ToLower(parts[0]) {
		case "temurin", "corretto", "zulu", "graalvm", "oracle", "microsoft", "liberica", "semeru", "dragonwell":
			return strings.ToLower(parts[0]), parts[1]
		}
	}
	return "temurin", input
}

// FormatVersionName returns directory name like "temurin-21.0.2"
func FormatVersionName(vendor, ver string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(vendor), ver)
}

// LocateJavaHome finds the true JAVA_HOME directory inside an extracted or existing JDK.
// On macOS, it's typically <dir>/Contents/Home or <dir>/<nested>/Contents/Home.
// On Linux/Windows, it's the directory containing 'bin/java'.
func LocateJavaHome(dir string) (string, error) {
	// Direct check
	binJava := filepath.Join(dir, "bin", "java")
	if runtime.GOOS == "windows" {
		binJava = filepath.Join(dir, "bin", "java.exe")
	}
	if _, err := os.Stat(binJava); err == nil {
		return dir, nil
	}

	// macOS Contents/Home
	macHome := filepath.Join(dir, "Contents", "Home")
	if _, err := os.Stat(filepath.Join(macHome, "bin", "java")); err == nil {
		return macHome, nil
	}

	var foundHome string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(dir, path)
		if strings.Count(rel, string(filepath.Separator)) > 3 {
			return filepath.SkipDir
		}

		if strings.HasSuffix(path, filepath.Join("Contents", "Home")) {
			testBin := filepath.Join(path, "bin", "java")
			if _, e := os.Stat(testBin); e == nil {
				foundHome = path
				return filepath.SkipAll
			}
		}

		testBin := filepath.Join(path, "bin", "java")
		if runtime.GOOS == "windows" {
			testBin = filepath.Join(path, "bin", "java.exe")
		}
		if _, e := os.Stat(testBin); e == nil {
			foundHome = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return "", err
	}

	if foundHome != "" {
		return foundHome, nil
	}

	return "", fmt.Errorf("could not locate Java binaries in %s", dir)
}

// JDKInfo stores parsed release metadata from a JDK
type JDKInfo struct {
	JavaHome    string
	Version     string
	Vendor      string
	SuggestedID string
}

// InspectJDK discovers JavaHome and parses release metadata or runs java -version
func InspectJDK(dir string) (*JDKInfo, error) {
	javaHome, err := LocateJavaHome(dir)
	if err != nil {
		return nil, err
	}

	info := &JDKInfo{
		JavaHome: javaHome,
		Vendor:   "custom",
		Version:  "unknown",
	}

	// 1. Try reading the 'release' file inside javaHome (fast & reliable without running java)
	releaseFile := filepath.Join(javaHome, "release")
	if data, err := os.ReadFile(releaseFile); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "JAVA_VERSION=") {
				ver := strings.Trim(strings.TrimPrefix(line, "JAVA_VERSION="), "\"")
				info.Version = ver
			} else if strings.HasPrefix(line, "IMPLEMENTOR=") {
				vend := strings.ToLower(strings.Trim(strings.TrimPrefix(line, "IMPLEMENTOR="), "\""))
				if strings.Contains(vend, "eclipse") || strings.Contains(vend, "adoptium") {
					info.Vendor = "temurin"
				} else if strings.Contains(vend, "amazon") {
					info.Vendor = "corretto"
				} else if strings.Contains(vend, "azul") {
					info.Vendor = "zulu"
				} else if strings.Contains(vend, "homebrew") {
					info.Vendor = "homebrew"
				} else if strings.Contains(vend, "oracle") {
					info.Vendor = "oracle"
				} else if strings.Contains(vend, "graalvm") {
					info.Vendor = "graalvm"
				}
			}
		}
	}

	// 2. Fall back to executing <javaHome>/bin/java -version
	if info.Version == "unknown" {
		binJava := filepath.Join(javaHome, "bin", "java")
		if runtime.GOOS == "windows" {
			binJava = filepath.Join(javaHome, "bin", "java.exe")
		}
		cmd := exec.Command(binJava, "-version")
		out, err := cmd.CombinedOutput()
		if err == nil {
			// e.g. openjdk version "21.0.2" ...
			re := regexp.MustCompile(`version\s+"([^"]+)"`)
			matches := re.FindStringSubmatch(string(out))
			if len(matches) > 1 {
				info.Version = matches[1]
			}
		}
	}

	if info.Version != "unknown" {
		info.SuggestedID = fmt.Sprintf("%s-%s", info.Vendor, info.Version)
	} else {
		info.SuggestedID = filepath.Base(dir)
	}

	return info, nil
}
