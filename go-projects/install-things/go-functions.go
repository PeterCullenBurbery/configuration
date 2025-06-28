package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"fmt"
	"path/filepath"
	"archive/zip"
	
	yekazip "github.com/yeka/zip"
	gofunctions "github.com/PeterCullenBurbery/go-functions"
)

// --- Helper functions ---

func getCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if result, ok := v.(map[string]interface{}); ok {
				return result
			}
		}
	}
	return nil
}

func getCaseInsensitiveList(m map[string]interface{}, key string) []string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			raw, ok := v.([]interface{})
			if !ok {
				return nil
			}
			var result []string
			for _, val := range raw {
				if s, ok := val.(string); ok {
					result = append(result, s)
				}
			}
			return result
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

func getNestedString(m map[string]interface{}, key string) string {
	if val := getCaseInsensitiveString(m, key); val != "" {
		return val
	}
	if sub := getCaseInsensitiveMap(m, key); sub != nil {
		for _, v := range sub {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	return getCaseInsensitiveMap(m, key)
}

func downloadFile(dest, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func formatTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// New helper to run each individual script:
func runPowerShellScript(filename, content string, logFile *os.File) {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("❌ Failed to write script %s: %v", filename, err)
	}
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", filename)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to run %s: %v", filename, err)
	}
	log.Printf("✅ Finished: %s", filename)
}

func unzipWithPassword(zipPath, destDir, password string) error {
	r, err := yekazip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		f.SetPassword(password)
		outPath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(outPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", outPath)
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(outPath, os.ModePerm)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		dstFile, _ := os.Create(outPath)
		defer dstFile.Close()
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleNirsoft(globalLogDir string, perAppLogs map[string]interface{}, globalDownloadDir string, perAppDownloads map[string]interface{}, modulePath string) {
	appKey := "nirsoft"
	subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

	// Explicit javac/java
	javac := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\javac.exe`
	java := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\java.exe`

	// Base dir (excluded from Defender)
	nirsoftBaseDir := filepath.Join(globalDownloadDir, subDownload)
	excludeFromDefender(nirsoftBaseDir)

	// Timestamped download subdir
	rawTimestampDownload, err := gofunctions.DateTimeStamp(javac, java)
	if err != nil {
		log.Fatalf("❌ Failed to generate download timestamp: %v", err)
	}
	downloadTimestamp := gofunctions.SafeTimeStamp(rawTimestampDownload, 1)
	nirsoftDownloadDir := filepath.Join(nirsoftBaseDir, downloadTimestamp)

	// ZIP metadata
	zipName := "nirsoft_package_enc_1.30.19.zip"
	zipPath := filepath.Join(nirsoftDownloadDir, zipName)
	zipURL := "https://github.com/PeterCullenBurbery/configuration/raw/main/host/password-protected/" + zipName
	password := "nirsoft9876$"

	// Timestamped extract dir
	rawTimestampExtract, err := gofunctions.DateTimeStamp(javac, java)
	if err != nil {
		log.Fatalf("❌ Failed to generate extraction timestamp: %v", err)
	}
	extractTimestamp := gofunctions.SafeTimeStamp(rawTimestampExtract, 1)
	extractDir := filepath.Join(nirsoftDownloadDir, extractTimestamp)

	// Optional log path
	var logPath string
	if subLog != "" && globalLogDir != "" {
		logDir := filepath.Join(globalLogDir, subLog)
		if err := os.MkdirAll(logDir, os.ModePerm); err == nil {
			logFileName := fmt.Sprintf("nirsoft_%s.log", downloadTimestamp)
			logPath = filepath.Join(logDir, logFileName)
		}
	}

	// Create folder for download
	if err := os.MkdirAll(nirsoftDownloadDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create download directory: %v", err)
	}
	log.Printf("📁 Creating download folder:\n↳ %s", nirsoftDownloadDir)

	// Download if needed
	if !fileExists(zipPath) {
		log.Printf("⬇️ Downloading: %s", zipURL)
		if err := downloadFile(zipPath, zipURL); err != nil {
			log.Fatalf("❌ Failed to download Nirsoft ZIP: %v", err)
		}
		log.Printf("✅ ZIP downloaded to: %s", zipPath)
	} else {
		log.Printf("📁 ZIP already exists: %s", zipPath)
	}

	// Extract
	if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create extract directory: %v", err)
	}
	log.Printf("📁 Creating extract folder:\n↳ %s", extractDir)

	if err := unzipWithPassword(zipPath, extractDir, password); err != nil {
		log.Fatalf("❌ Failed to extract password-protected ZIP: %v", err)
	}
	log.Println("✅ Extraction complete!")

	// Summary
	log.Printf("📦 Extracted Nirsoft package to:\n↳ %s", extractDir)
	if logPath != "" {
		log.Printf("📝 Nirsoft log path:\n↳ %s", logPath)
	}
}

// Exclude directory from Defender
func excludeFromDefender(path string) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("⚠️ Failed to exclude from Defender: %s\n↳ %v", path, err)
	} else {
		log.Printf("🛡️ Added Defender exclusion:\n↳ %s", path)
	}
}

func handleSQLDeveloper(globalLogDir string, perAppLogs map[string]interface{}, globalDownloadDir string, perAppDownloads map[string]interface{}, modulePath string) {
	appKey := "sql developer"

	subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

	timestamp := formatTimestamp()
	sqlDownloadDir := filepath.Join(globalDownloadDir, subDownload)

	// Optional log path
	var sqlLogPath string
	if subLog != "" && globalLogDir != "" {
		logDir := filepath.Join(globalLogDir, subLog)
		if err := os.MkdirAll(logDir, os.ModePerm); err == nil {
			logFileName := fmt.Sprintf("sqldeveloper_%s.log", timestamp)
			sqlLogPath = filepath.Join(logDir, logFileName)
		}
	}

	_ = os.MkdirAll(sqlDownloadDir, os.ModePerm)

	zipName := "sqldeveloper-24.3.1.347.1826-x64.zip"
	zipPath := filepath.Join(sqlDownloadDir, zipName)
	extractDir := filepath.Join(sqlDownloadDir, strings.TrimSuffix(zipName, filepath.Ext(zipName)))
	installerURL := "https://download.oracle.com/otn_software/java/sqldeveloper/" + zipName

	if !fileExists(zipPath) {
		log.Printf("🌐 Downloading SQL Developer from: %s", installerURL)
		if err := downloadFile(zipPath, installerURL); err != nil {
			log.Fatalf("❌ Download failed: %v", err)
		}
		log.Println("✅ Downloaded SQL Developer.")
	} else {
		log.Println("📁 SQL Developer ZIP already present.")
	}

	log.Printf("📦 Extracting SQL Developer to: %s", extractDir)
	if err := unzip(zipPath, extractDir); err != nil {
		log.Fatalf("❌ Extraction failed: %v", err)
	}
	log.Println("✅ SQL Developer extracted.")

	if sqlLogPath != "" {
		log.Printf("📝 SQL Developer log path: %s", sqlLogPath)
	}

	// Step: Create shortcut using the module's New-DesktopShortcut
	exePath := filepath.Join(extractDir, "sqldeveloper", "sqldeveloper.exe")
	psContent := fmt.Sprintf(`Import-Module '%s'
New-DesktopShortcut -TargetPath '%s' -Description 'Oracle SQL Developer'
`, modulePath, exePath)

	psScriptName := "create-sqldeveloper-shortcut.ps1"
	if err := os.WriteFile(psScriptName, []byte(psContent), 0644); err != nil {
		log.Fatalf("❌ Failed to write shortcut script: %v", err)
	}

	log.Println("📌 Creating SQL Developer desktop shortcut...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to create SQL Developer shortcut: %v", err)
	}
	log.Println("✅ SQL Developer shortcut created on desktop.")
}

func handleCherryTree(globalLogDir string, perAppLogs map[string]interface{}, globalDownloadDir string, perAppDownloads map[string]interface{}, modulePath string) {
	appKey := "cherry tree"
	subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))
	timestamp := formatTimestamp()

	logDir := filepath.Join(globalLogDir, subLog)
	logFileName := fmt.Sprintf("cherrytree_%s.log", timestamp)
	cherryLogPath := filepath.Join(logDir, logFileName)
	cherryInstallPath := filepath.Join(globalDownloadDir, subDownload)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create CherryTree log directory: %v", err)
	}
	if err := os.MkdirAll(cherryInstallPath, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create CherryTree install directory: %v", err)
	}

	installerPath := filepath.Join(cherryInstallPath, "cherrytree_1.5.0.0_win64_setup.exe")
	installerURL := "https://www.giuspen.net/software/cherrytree_1.5.0.0_win64_setup.exe"

	if !fileExists(installerPath) {
		log.Printf("🌐 Downloading CherryTree from: %s", installerURL)
		if err := downloadFile(installerPath, installerURL); err != nil {
			log.Fatalf("❌ Download failed: %v", err)
		}
		log.Println("✅ Downloaded CherryTree.")
	} else {
		log.Println("📁 CherryTree installer already present.")
	}

	log.Printf("📝 CherryTree log path: %s", cherryLogPath)
	psContent := fmt.Sprintf("Import-Module '%s'\nInstall-CherryTree -log '%s' -installPath '%s'\n", modulePath, cherryLogPath, cherryInstallPath)
	runPowerShellScript("install-cherrytree.ps1", psContent, nil)
}

func handleMiniconda(globalDownloadDir string, perAppDownloads map[string]interface{}, modulePath string) {
	appKey := "python"
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

	minicondaInstallPath := filepath.Join(globalDownloadDir, subDownload)
	installerPath := filepath.Join(minicondaInstallPath, "Miniconda3-latest-Windows-x86_64.exe")
	installerURL := "https://repo.anaconda.com/miniconda/Miniconda3-latest-Windows-x86_64.exe"

	if err := os.MkdirAll(minicondaInstallPath, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create Miniconda install path: %v", err)
	}

	if !fileExists(installerPath) {
		log.Printf("🌐 Downloading Miniconda from: %s", installerURL)
		if err := downloadFile(installerPath, installerURL); err != nil {
			log.Fatalf("❌ Download failed: %v", err)
		}
		log.Println("✅ Downloaded Miniconda.")
	} else {
		log.Println("📁 Miniconda installer already present.")
	}

	psContent := fmt.Sprintf(`Import-Module '%s'

Install-Miniconda -InstallerPath '%s'

# Add Python and Pip to PATH
Add-ToPath -PathToAdd 'C:\ProgramData\Miniconda3\python.exe'
Add-ToPath -PathToAdd 'C:\ProgramData\Miniconda3\Scripts\pip3.exe'
`, modulePath, installerPath)

	runPowerShellScript("install-miniconda.ps1", psContent, nil)
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}