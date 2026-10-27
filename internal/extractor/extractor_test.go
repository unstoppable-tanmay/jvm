package extractor_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/unstoppable-tanmay/jvm/internal/extractor"
)

func TestExtractTarGz(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_targz_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy tar.gz in memory
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	content := "mock java binary content"
	hdr := &tar.Header{
		Name: "test-jdk/bin/java",
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gw.Close()

	tarPath := filepath.Join(tempDir, "sample.tar.gz")
	if err := os.WriteFile(tarPath, buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tempDir, "extracted")
	if err := extractor.Extract(tarPath, destDir); err != nil {
		t.Fatalf("Extract tar.gz failed: %v", err)
	}

	extractedFile := filepath.Join(destDir, "test-jdk", "bin", "java")
	data, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(data) != content {
		t.Errorf("content mismatch, expected %q, got %q", content, string(data))
	}
}

func TestExtractZip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_zip_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy zip in memory
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	content := "mock windows java binary"
	w, err := zw.Create("jdk-windows/bin/java.exe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	zw.Close()

	zipPath := filepath.Join(tempDir, "sample.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tempDir, "extracted")
	if err := extractor.Extract(zipPath, destDir); err != nil {
		t.Fatalf("Extract zip failed: %v", err)
	}

	extractedFile := filepath.Join(destDir, "jdk-windows", "bin", "java.exe")
	data, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted zip file: %v", err)
	}
	if string(data) != content {
		t.Errorf("content mismatch, expected %q, got %q", content, string(data))
	}
}
