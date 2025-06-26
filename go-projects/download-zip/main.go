package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Replace the "blob" URL with the "raw" GitHub URL
	url := "https://raw.githubusercontent.com/PeterCullenBurbery/configuration/main/go-projects/zip/nirsoft_package_enc_1.30.19.zip"

	// Set download path
	destDir := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\go-projects\download-zip\test\folder-001`
	destFile := "nirsoft_package_enc_1.30.19.zip"
	fullPath := filepath.Join(destDir, destFile)

	// Ensure directory exists
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// Fetch the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Failed to download: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Bad status: %s\n", resp.Status)
		os.Exit(1)
	}

	// Write content to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to save file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Downloaded to %s\n", fullPath)
}