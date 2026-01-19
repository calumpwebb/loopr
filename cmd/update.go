package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/calumpwebb/loopr/internal/ui"
)

const githubRepo = "calumpwebb/loopr"

// Update checks for and installs the latest version of loopr
func Update(currentVersion string) error {
	// 1. Get current binary path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// 2. Fetch latest release from GitHub API
	fmt.Println()
	fmt.Println("Checking for updates...")
	latestVersion, downloadURL, err := fetchLatestRelease()
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Failed to check for updates"))
		return err
	}

	// 3. Compare versions
	if currentVersion == latestVersion {
		fmt.Println(ui.SuccessStyle.Render("✓ Already on latest version: " + currentVersion))
		return nil
	}

	// 4. Show update prompt
	fmt.Println()
	fmt.Printf("Update available: %s → %s\n", currentVersion, latestVersion)
	fmt.Println()
	if !ui.PromptUpdate() {
		fmt.Println("Cancelled.")
		return nil
	}

	// 5. Download new binary to temp file
	fmt.Println()
	fmt.Println("Downloading...")
	tmpFile, err := downloadBinary(downloadURL)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Download failed"))
		return err
	}
	defer os.Remove(tmpFile)

	// 6. Replace current binary (atomic if possible)
	fmt.Println("Installing...")
	if err := replaceBinary(execPath, tmpFile); err != nil {
		fmt.Println(ui.ErrorStyle.Render("✗ Installation failed"))
		return err
	}

	// 7. Success message
	fmt.Println(ui.SuccessStyle.Render("✓ Updated to " + latestVersion))
	fmt.Println()
	fmt.Println("Run 'loopr version' to verify")
	fmt.Println()

	return nil
}

// fetchLatestRelease fetches the latest release info from GitHub API
func fetchLatestRelease() (version, downloadURL string, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", fmt.Errorf("failed to parse release info: %w", err)
	}

	// Find matching binary for current platform
	binaryName := fmt.Sprintf("loopr-%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			return release.TagName, asset.BrowserDownloadURL, nil
		}
	}

	return "", "", fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
}

// downloadBinary downloads a binary from a URL to a temporary file
func downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "loopr-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to save binary: %w", err)
	}

	// Make executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	return tmpFile.Name(), nil
}

// replaceBinary replaces the target binary with the source binary atomically
func replaceBinary(target, source string) error {
	// Create backup
	backup := target + ".backup"
	if err := os.Rename(target, backup); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Move new binary to target location
	if err := os.Rename(source, target); err != nil {
		// Rollback on failure
		os.Rename(backup, target)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Remove backup on success
	os.Remove(backup)
	return nil
}
