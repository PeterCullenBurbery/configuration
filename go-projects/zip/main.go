package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	sourceDir := `C:\Users\Administrator\Desktop\testing\zip\folder-001`
	outputZip := `C:\Users\Administrator\Desktop\go-projects\zip\nirsoft_package_enc_1.30.19.zip`

	err := zipFolder(sourceDir, outputZip)
	if err != nil {
		fmt.Printf("❌ Failed to zip folder: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Folder zipped successfully.")
}

func zipFolder(sourceDir, zipPath string) error {
	// Create the zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("create zip: %w", err)
	}
	defer zipFile.Close()

	// Initialize the zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("⚠️ Skipping %s due to error: %v\n", path, err)
			return nil
		}

		// Skip directories, they are created implicitly
		if info.IsDir() {
			return nil
		}

		// Open file
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("⚠️ Could not open file %s: %v\n", path, err)
			return nil // Skip unreadable files
		}
		defer file.Close()

		// Create zip header
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("create file header: %w", err)
		}
		header.Name = filepath.ToSlash(relPath)
		header.Method = zip.Deflate

		// Create writer for the file in zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("create zip entry: %w", err)
		}

		// Copy file content to zip
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("copy to zip: %w", err)
		}

		return nil
	})
}
