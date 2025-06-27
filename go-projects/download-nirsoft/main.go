package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PeterCullenBurbery/go-functions" // Adjust if using local path
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if submap, ok := v.(map[string]interface{}); ok {
				return submap
			}
		}
	}
	return nil
}

func getCaseInsensitiveString(m map[string]interface{}, key string) string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}

func getNestedString(m map[string]interface{}, parentKey string) string {
	for k, v := range m {
		if strings.EqualFold(k, parentKey) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}

func getNestedMap(m map[string]interface{}, parentKey string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, parentKey) {
			if submap, ok := v.(map[string]interface{}); ok {
				return submap
			}
		}
	}
	return nil
}

func downloadFile(dst string, url string) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Invoke-WebRequest -Uri \"%s\" -OutFile \"%s\"", url, dst))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func unzip(src, dest string) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Expand-Archive -Path \"%s\" -DestinationPath \"%s\" -Force", src, dest))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	installPath := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\install.yaml`

	data, err := os.ReadFile(installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read install.yaml: %v", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	installSection := getCaseInsensitiveMap(raw, "install")
	if installSection == nil {
		log.Fatal("❌ Missing 'install' section.")
	}

	logs := getCaseInsensitiveMap(installSection, "logs")
	downloads := getCaseInsensitiveMap(installSection, "downloads")

	globalLogDir := strings.TrimSpace(getNestedString(logs, "global log directory"))
	globalDownloadDir := strings.TrimSpace(getNestedString(downloads, "global download directory"))
	perAppLogs := getNestedMap(logs, "per app log directories")
	perAppDownloads := getNestedMap(downloads, "per app download directories")

	appKey := "nirsoft"
	subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

	javac := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\javac.exe`
	java := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\java.exe`

	// Generate timestamps
	rawDownload, err := gofunctions.DateTimeStamp(javac, java)
	if err != nil {
		log.Fatalf("❌ Timestamp (download) error: %v", err)
	}
	downloadTimestamp := gofunctions.SafeTimeStamp(rawDownload, 1)

	rawExtract, err := gofunctions.DateTimeStamp(javac, java)
	if err != nil {
		log.Fatalf("❌ Timestamp (extract) error: %v", err)
	}
	extractTimestamp := gofunctions.SafeTimeStamp(rawExtract, 1)

	// Construct paths
	nirsoftDownloadDir := filepath.Join(globalDownloadDir, subDownload, downloadTimestamp)
	zipName := "nirsoft_package_enc_1.30.19.zip"
	zipPath := filepath.Join(nirsoftDownloadDir, zipName)
	extractDir := filepath.Join(nirsoftDownloadDir, extractTimestamp)

	_ = os.MkdirAll(nirsoftDownloadDir, os.ModePerm)

	logFilePath := ""
	if globalLogDir != "" && subLog != "" {
		logDir := filepath.Join(globalLogDir, subLog)
		_ = os.MkdirAll(logDir, os.ModePerm)
		logFilePath = filepath.Join(logDir, fmt.Sprintf("nirsoft_%s.log", downloadTimestamp))

		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			defer logFile.Close()
			log.SetOutput(io.MultiWriter(os.Stdout, logFile))
		}
	}

	log.Printf("📁 Creating download folder:\n↳ %s", nirsoftDownloadDir)

	if !fileExists(zipPath) {
		log.Printf("⬇️ Downloading: %s", "https://github.com/PeterCullenBurbery/configuration/raw/main/host/"+zipName)
		if err := downloadFile(zipPath, "https://github.com/PeterCullenBurbery/configuration/raw/main/host/"+zipName); err != nil {
			log.Fatalf("❌ Download failed: %v", err)
		}
		log.Printf("✅ ZIP downloaded to: %s", zipPath)
	} else {
		log.Printf("📁 ZIP already exists: %s", zipPath)
	}

	log.Printf("📁 Creating extract folder:\n↳ %s", extractDir)
	if err := unzip(zipPath, extractDir); err != nil {
		log.Fatalf("❌ Extract failed: %v", err)
	}
	log.Println("✅ Extraction complete!")
	log.Printf("📦 Extracted to: %s", extractDir)
	if logFilePath != "" {
		log.Printf("📝 Log saved at: %s", logFilePath)
	}
}