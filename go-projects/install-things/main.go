package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		if funcName != "" {
			log.Printf("➡️  Starting: %s → %s", label, funcName)
		} else if funcName == "" {
			log.Printf("➡️  Starting: %s", label)
		}

		switch {
		case strings.EqualFold(label, "SQL Developer"):
			handleSQLDeveloper(globalLogDir, perAppLogs, globalDownloadDir, perAppDownloads, *modulePath)

		case strings.EqualFold(label, "Nirsoft"):
			handleNirsoft(globalLogDir, perAppLogs, globalDownloadDir, perAppDownloads, *modulePath)

		case strings.EqualFold(funcName, "Install-CherryTree"):
			handleCherryTree(globalLogDir, perAppLogs, globalDownloadDir, perAppDownloads, *modulePath)

			// Miniconda does not take logs.
		case strings.EqualFold(funcName, "Install-Miniconda"):
			handleMiniconda(globalDownloadDir, perAppDownloads, *modulePath)
		default:
			scriptName := fmt.Sprintf("install-%s.ps1", strings.ToLower(strings.ReplaceAll(label, " ", "-")))
			psContent := fmt.Sprintf("Import-Module '%s'\n%s\n", *modulePath, funcName)
			runPowerShellScript(scriptName, psContent, logFile)
		}
	}

	log.Println("🎉 All installations completed.")
}

// --- Helper functions ---

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
