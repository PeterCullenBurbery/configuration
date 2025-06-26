package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Direct download URL
	url := "https://raw.githubusercontent.com/PeterCullenBurbery/configuration/main/go-projects/zip/nirsoft_package_enc_1.30.19.zip"

	// Paths
	downloadDir := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\go-projects\download-zip\test\folder-003`
	extractDir := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\go-projects\download-zip\test\folder-003`
	zipFileName := "nirsoft_package_enc_1.30.19.zip"
	zipFilePath := filepath.Join(downloadDir, zipFileName)

	// Ensure directories exist
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		fmt.Printf("❌ Failed to create download directory: %v\n", err)
		return
	}
	if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
		fmt.Printf("❌ Failed to create extract directory: %v\n", err)
		return
	}

	// Download ZIP
	fmt.Println("🌐 Downloading ZIP file...")
	out, err := os.Create(zipFilePath)
	if err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		return
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Failed to download file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Bad HTTP response: %s\n", resp.Status)
		return
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Printf("❌ Failed to save file: %v\n", err)
		return
	}
	fmt.Printf("✅ Downloaded to: %s\n", zipFilePath)

	// Extract ZIP
	fmt.Println("📦 Extracting ZIP file...")
	if err := extractZip(zipFilePath, extractDir); err != nil {
		fmt.Printf("❌ Extraction failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Extracted to: %s\n", extractDir)
}

// extractZip unzips a .zip archive to the target directory
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, inFile); err != nil {
			return err
		}
	}
	return nil
}