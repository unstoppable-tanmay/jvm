package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// Download downloads a remote file to targetPath while displaying a progress bar, and verifies SHA256 if provided
func Download(url, targetPath, expectedChecksum string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "jvm-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to initiate download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from server: %s", resp.Status)
	}

	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	hasher := sha256.New()
	multiWriter := io.MultiWriter(out, bar, hasher)

	if _, err := io.Copy(multiWriter, resp.Body); err != nil {
		return fmt.Errorf("download error: %w", err)
	}

	fmt.Println() // Newline after progress bar finishes

	if expectedChecksum != "" {
		actualChecksum := hex.EncodeToString(hasher.Sum(nil))
		if actualChecksum != expectedChecksum {
			os.Remove(targetPath)
			return fmt.Errorf("checksum mismatch!\nexpected: %s\nactual:   %s", expectedChecksum, actualChecksum)
		}
	}

	return nil
}
