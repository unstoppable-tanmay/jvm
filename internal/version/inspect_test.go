package version_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unstoppable-tanmay/jvm/internal/version"
)

func TestInspectJDK_StandardLayout(t *testing.T) {
	tempJDK, err := os.MkdirTemp("", "mock_jdk_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempJDK)

	binDir := filepath.Join(tempJDK, "bin")
	_ = os.MkdirAll(binDir, 0755)
	_ = os.WriteFile(filepath.Join(binDir, "java"), []byte("#!/bin/sh\n"), 0755)
	_ = os.WriteFile(filepath.Join(binDir, "java.exe"), []byte("mock binary"), 0755)

	releaseContent := `JAVA_VERSION="17.0.12"
IMPLEMENTOR="Eclipse Adoptium"
`
	_ = os.WriteFile(filepath.Join(tempJDK, "release"), []byte(releaseContent), 0644)

	info, err := version.InspectJDK(tempJDK)
	if err != nil {
		t.Fatalf("InspectJDK failed: %v", err)
	}

	if info.Vendor != "temurin" {
		t.Errorf("expected vendor temurin, got %s", info.Vendor)
	}
	if info.Version != "17.0.12" {
		t.Errorf("expected version 17.0.12, got %s", info.Version)
	}
	if info.SuggestedID != "temurin-17.0.12" {
		t.Errorf("expected suggested ID temurin-17.0.12, got %s", info.SuggestedID)
	}
}

func TestInspectJDK_MacOSBundleLayout(t *testing.T) {
	tempBundle, err := os.MkdirTemp("", "mock_bundle_*.jdk")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempBundle)

	homeDir := filepath.Join(tempBundle, "Contents", "Home")
	binDir := filepath.Join(homeDir, "bin")
	_ = os.MkdirAll(binDir, 0755)
	_ = os.WriteFile(filepath.Join(binDir, "java"), []byte("#!/bin/sh\n"), 0755)
	_ = os.WriteFile(filepath.Join(binDir, "java.exe"), []byte("mock binary"), 0755)

	releaseContent := `JAVA_VERSION="21.0.2"
IMPLEMENTOR="Amazon.com Inc."
`
	_ = os.WriteFile(filepath.Join(homeDir, "release"), []byte(releaseContent), 0644)

	info, err := version.InspectJDK(tempBundle)
	if err != nil {
		t.Fatalf("InspectJDK failed: %v", err)
	}

	if info.JavaHome != homeDir {
		t.Errorf("expected JavaHome %s, got %s", homeDir, info.JavaHome)
	}
	if info.Vendor != "corretto" {
		t.Errorf("expected vendor corretto, got %s", info.Vendor)
	}
	if info.Version != "21.0.2" {
		t.Errorf("expected version 21.0.2, got %s", info.Version)
	}
}
