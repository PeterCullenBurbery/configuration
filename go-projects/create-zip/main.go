package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"bytes"

	"github.com/PeterCullenBurbery/go-functions"
	"gopkg.in/yaml.v3"
)

func main() {
	yamlPath := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\install.yaml`
	zipSource := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\host\nirsoft_package_enc_1.30.19.zip`

	// Load YAML
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}
	var root map[string]interface{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	// Extract relevant fields from YAML
	install := gofunctions.GetCaseInsensitiveMap(root, "install")
	downloads := gofunctions.GetCaseInsensitiveMap(install, "downloads")
	globalDownloadDir := strings.TrimSpace(gofunctions.GetCaseInsensitiveString(downloads, "global download directory"))
	perAppDownloads := gofunctions.GetCaseInsensitiveMap(downloads, "per app download directories")
	subDownload := strings.TrimSpace(gofunctions.GetCaseInsensitiveString(perAppDownloads, "Nirsoft"))

	// Add Defender exclusion for C:\downloads\nirsoft
	nirsoftBaseDir := filepath.Join(globalDownloadDir, subDownload)
	if err := addDefenderExclusion(nirsoftBaseDir); err != nil {
		log.Printf("⚠️ Failed to add Defender exclusion for %s: %v", nirsoftBaseDir, err)
	} else {
		fmt.Println("✅ Added Defender exclusion for:", nirsoftBaseDir)
	}

	// First SafeTimeStamp (for copy)
	rawTimestampCopy, err := gofunctions.DateTimeStamp()
	if err != nil {
		log.Fatalf("❌ Failed to generate copy timestamp: %v", err)
	}
	safeTimestampCopy := gofunctions.SafeTimeStamp(rawTimestampCopy, 1)

	// Prepare copy target
	copyDir := filepath.Join(nirsoftBaseDir, safeTimestampCopy)
	zipCopyPath := filepath.Join(copyDir, filepath.Base(zipSource))

	if err := os.MkdirAll(copyDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create copy directory: %v", err)
	}

	// Copy the ZIP
	if err := copyFile(zipSource, zipCopyPath); err != nil {
		log.Fatalf("❌ Failed to copy ZIP file: %v", err)
	}
	fmt.Println("✅ ZIP copied to:", zipCopyPath)

	// Second SafeTimeStamp (for extract)
	rawTimestampExtract, err := gofunctions.DateTimeStamp()
	if err != nil {
		log.Fatalf("❌ Failed to generate extract timestamp: %v", err)
	}
	safeTimestampExtract := gofunctions.SafeTimeStamp(rawTimestampExtract, 1)

	// Extract to: copyDir/safeTimestampExtract
	extractDir := filepath.Join(copyDir, safeTimestampExtract)
	if err := unzip(zipCopyPath, extractDir); err != nil {
		log.Fatalf("❌ Failed to extract ZIP: %v", err)
	}
	fmt.Println("✅ ZIP extracted to:", extractDir)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// unzip extracts zipPath into destDir
func unzip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		fullPath := filepath.Join(destDir, file.Name)
		if !strings.HasPrefix(fullPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fullPath)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}

func addDefenderExclusion(path string) error {
	checkCmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-Command", fmt.Sprintf(`(Get-MpPreference).ExclusionPath -contains "%s"`, path))

	var out bytes.Buffer
	checkCmd.Stdout = &out
	checkCmd.Stderr = &out
	if err := checkCmd.Run(); err == nil && strings.Contains(out.String(), "True") {
		log.Printf("ℹ️ Defender exclusion already exists:\n↳ %s", path)
		return nil
	}

	psCommand := fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, path)
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}