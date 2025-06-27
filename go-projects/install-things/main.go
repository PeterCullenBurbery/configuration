package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gofunctions "github.com/PeterCullenBurbery/go-functions"
	"gopkg.in/yaml.v3"
)

func main() {
	installPath := flag.String("install", "", "Path to install.yaml (required)")
	modulePath := flag.String("module", "", "Path to PowerShell module (.psm1) (required)")
	logPath := flag.String("log", "", "Path to runtime execution log file (required)")
	flag.Parse()

	if *installPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ --install, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	log.Println("🍫 Installing Chocolatey using Install-Choco...")
	psContent := fmt.Sprintf("Import-Module '%s'\nInstall-Choco\n", *modulePath)
	runPowerShellScript("install-choco.ps1", psContent, logFile)
	log.Println("✅ Chocolatey installation complete.")

	var raw map[string]interface{}
	data, err := os.ReadFile(*installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	installSection := getCaseInsensitiveMap(raw, "install")
	if installSection == nil {
		log.Fatal("❌ Missing 'install' section.")
	}

	programs := getCaseInsensitiveList(installSection, "programs to install")
	logs := getCaseInsensitiveMap(installSection, "logs")
	downloads := getCaseInsensitiveMap(installSection, "downloads")

	globalLogDir := strings.TrimSpace(getNestedString(logs, "global log directory"))
	perAppLogs := getNestedMap(logs, "per app log directories")
	globalDownloadDir := strings.TrimSpace(getNestedString(downloads, "global download directory"))
	perAppDownloads := getNestedMap(downloads, "per app download directories")

	if globalLogDir == "" || globalDownloadDir == "" {
		log.Fatal("❌ Missing 'global log directory' or 'global download directory'.")
	}

	_ = os.MkdirAll(globalLogDir, os.ModePerm)
	_ = os.MkdirAll(globalDownloadDir, os.ModePerm)

	for _, label := range programs {
		if strings.EqualFold(label, "Choco") {
			continue // ✅ Already handled above
		}

		funcName := toInstallFunctionName(label)
		log.Printf("➡️  Starting: %s → %s", label, funcName)

		switch {
		case strings.EqualFold(label, "SQL Developer"):
			handleSQLDeveloper(globalLogDir, perAppLogs, globalDownloadDir, perAppDownloads, *modulePath)

		case strings.EqualFold(label, "Nirsoft"):
			handleNirsoft(globalLogDir, perAppLogs, globalDownloadDir, perAppDownloads, *modulePath)

		case strings.EqualFold(funcName, "Install-CherryTree"):
			appKey := "cherry tree"
			subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
			subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))
			timestamp := formatTimestamp()

			logDir := filepath.Join(globalLogDir, subLog)
			logFileName := fmt.Sprintf("cherrytree_%s.log", timestamp)
			cherryLogPath := filepath.Join(logDir, logFileName)
			cherryInstallPath := filepath.Join(globalDownloadDir, subDownload)

			_ = os.MkdirAll(logDir, os.ModePerm)
			_ = os.MkdirAll(cherryInstallPath, os.ModePerm)

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
			psContent := fmt.Sprintf("Import-Module '%s'\nInstall-CherryTree -log '%s' -installPath '%s'\n", *modulePath, cherryLogPath, cherryInstallPath)
			runPowerShellScript("install-cherrytree.ps1", psContent, logFile)

		case strings.EqualFold(funcName, "Install-Miniconda"):
			appKey := "python"
			subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))
			minicondaInstallPath := filepath.Join(globalDownloadDir, subDownload)
			installerPath := filepath.Join(minicondaInstallPath, "Miniconda3-latest-Windows-x86_64.exe")
			installerURL := "https://repo.anaconda.com/miniconda/Miniconda3-latest-Windows-x86_64.exe"

			_ = os.MkdirAll(minicondaInstallPath, os.ModePerm)

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
`, *modulePath, installerPath)

			runPowerShellScript("install-miniconda.ps1", psContent, logFile)

		default:
			scriptName := fmt.Sprintf("install-%s.ps1", strings.ToLower(strings.ReplaceAll(label, " ", "-")))
			psContent := fmt.Sprintf("Import-Module '%s'\n%s\n", *modulePath, funcName)
			runPowerShellScript(scriptName, psContent, logFile)
		}
	}

	log.Println("🎉 All installations completed.")
}

// --- Helper functions ---

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

func toInstallFunctionName(label string) string {
	l := strings.ToLower(strings.NewReplacer(" ", "", "-", "", "+", "").Replace(label))
	switch l {
	case "powershell7":
		return "Install-PowerShell-7"
	case "vscode":
		return "Install-VSCode"
	case "7zip":
		return "Install-7Zip"
	case "voidtoolseverything":
		return "Install-Voidtools-Everything"
	case "winscp":
		return "Install-WinSCP"
	case "mobaxterm":
		return "Install-MobaXterm"
	case "choco", "chocolatey":
		return "Install-Choco"
	case "cherrytree":
		return "Install-CherryTree"
	case "go":
		return "Install-Go"
	case "notepadpp", "notepadplusplus", "notepad":
		return "Install-NotepadPP"
	case "sqlitebrowser", "sqlite", "sqlitebrowserforsqlite", "dbbrowser":
		return "Install-SQLiteBrowser"
	case "python", "miniconda":
		return "Install-Miniconda"
	case "java":
		return "Install-Java"
	default:
		return "Install-" + strings.Title(l)
	}
}

func handleNirsoft(globalLogDir string, perAppLogs map[string]interface{}, globalDownloadDir string, perAppDownloads map[string]interface{}, modulePath string) {
	appKey := "nirsoft"

	subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
	subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

	// Explicit javac and java paths
	javac := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\javac.exe`
	java := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin\java.exe`

	// Generate safe timestamp for download phase
	rawTimestampDownload, err := gofunctions.DateTimeStamp(javac, java)
	if err != nil {
		log.Fatalf("❌ Failed to generate download timestamp: %v", err)
	}
	downloadTimestamp := gofunctions.SafeTimeStamp(rawTimestampDownload, 1)

	// Paths
	nirsoftDownloadDir := filepath.Join(globalDownloadDir, subDownload, downloadTimestamp)
	zipName := "nirsoft_package_enc_1.30.19.zip"
	zipPath := filepath.Join(nirsoftDownloadDir, zipName)
	zipURL := "https://github.com/PeterCullenBurbery/configuration/raw/main/host/" + zipName

	// Extraction timestamp
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

	// Create download folder
	if err := os.MkdirAll(nirsoftDownloadDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create download directory: %v", err)
	}
	log.Printf("📁 Creating download folder:\n↳ %s", nirsoftDownloadDir)

	// Defender exclusion on download dir
	excludeFromDefender(nirsoftDownloadDir)

	// Download
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

	if err := unzip(zipPath, extractDir); err != nil {
		log.Fatalf("❌ Failed to extract ZIP: %v", err)
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
