package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PeterCullenBurbery/go-functions"
)

func main() {
	// Step 1: Generate timestamped folder using SafeTimeStamp
	rawDownloadStamp, err := gofunctions.DateTimeStamp()
	if err != nil {
		fmt.Println("❌ Failed to generate download timestamp:", err)
		return
	}
	downloadFolder := strings.TrimSpace(gofunctions.SafeTimeStamp(rawDownloadStamp, 1))

	baseDir := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\go-projects\download-zip\download`
	fullDownloadPath := filepath.Join(baseDir, downloadFolder)

	fmt.Println("📁 Creating download folder:")
	fmt.Println("↳", fullDownloadPath)
	fmt.Println("🔢 Length:", len(fullDownloadPath))

	if err := os.MkdirAll(fullDownloadPath, os.ModePerm); err != nil {
		fmt.Println("❌ Failed to create download folder:", err)
		return
	}

	// Step 2: Download the ZIP file
	url := "https://github.com/PeterCullenBurbery/configuration/raw/main/host/nirsoft_package_enc_1.30.19.zip"
	zipFileName := "nirsoft_package_enc_1.30.19.zip"
	zipPath := filepath.Join(fullDownloadPath, zipFileName)

	fmt.Println("⬇️ Downloading:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("❌ Failed to download:", err)
		return
	}
	defer resp.Body.Close()

	outFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Println("❌ Failed to create zip file:", err)
		return
	}
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		outFile.Close()
		fmt.Println("❌ Failed to write zip file:", err)
		return
	}
	outFile.Close()
	fmt.Println("✅ ZIP downloaded to:", zipPath)

	// Step 3: Create extraction folder (based on new timestamp)
	rawExtractStamp, err := gofunctions.DateTimeStamp()
	if err != nil {
		fmt.Println("❌ Failed to generate extract timestamp:", err)
		return
	}
	extractFolder := strings.TrimSpace(gofunctions.SafeTimeStamp(rawExtractStamp, 1))
	fullExtractPath := filepath.Join(fullDownloadPath, extractFolder)

	fmt.Println("📁 Creating extract folder:")
	fmt.Println("↳", fullExtractPath)
	fmt.Println("🔢 Length:", len(fullExtractPath))

	if err := os.MkdirAll(fullExtractPath, os.ModePerm); err != nil {
		fmt.Println("❌ Failed to create extraction folder:", err)
		return
	}

	// Step 4: Extract ZIP
	fmt.Println("📦 Extracting ZIP to:", fullExtractPath)
	if err := unzip(zipPath, fullExtractPath); err != nil {
		fmt.Println("❌ Failed to extract zip:", err)
		return
	}

	fmt.Println("✅ Extraction complete!")
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Create directory if necessary
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the file
		rc, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(outFile, rc); err != nil {
			rc.Close()
			outFile.Close()
			return err
		}
		rc.Close()
		outFile.Close()
	}
	return nil
}