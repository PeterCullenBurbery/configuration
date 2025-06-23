package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	// --- CLI Arguments ---
	installPath := flag.String("install", "", "Path to install.yaml (required)")
	modulePath := flag.String("module", "", "Path to PowerShell module (.psm1) (required)")
	logPath := flag.String("log", "", "Path to log file (required)")
	flag.Parse()

	if *installPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ --install, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

	// --- Log Setup ---
	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// --- Load YAML ---
	log.Println("📦 Reading:", *installPath)
	data, err := os.ReadFile(*installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read install.yaml: %v", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Fatalf("❌ Failed to parse install.yaml: %v", err)
	}

	// --- Extract 'install' section ---
	installSection := getCaseInsensitiveMap(raw, "install")
	if installSection == nil {
		log.Fatalf("❌ 'install' section not found in YAML.")
	}

	// --- Extract 'programs to install' as ordered list ---
	programsList := getCaseInsensitiveList(installSection, "programs to install")
	if programsList == nil {
		log.Fatalf("❌ 'programs to install' section not found or is not a list.")
	}

	// --- Extract optional 'logs' and 'downloads' ---
	logsPath := getCaseInsensitiveString(installSection, "logs")
	downloadsPath := getCaseInsensitiveString(installSection, "downloads")

	log.Printf("🪵 Logs path: %s", logsPath)
	log.Printf("📥 Downloads path: %s", downloadsPath)

	// --- Process Installers in Order ---
	var psFunctions []string
	for _, rawKey := range programsList {
		funcName := toInstallFunctionName(rawKey)
		psFunctions = append(psFunctions, funcName)
		log.Printf("✔️  Queued installer: %s → %s", rawKey, funcName)
	}

	if len(psFunctions) == 0 {
		log.Println("⚠️ No programs listed. Exiting.")
		return
	}

	// --- Generate PowerShell script ---
	psScript := fmt.Sprintf("Import-Module '%s'\n", *modulePath)
	for _, fn := range psFunctions {
		psScript += fn + "\n"
	}

	tempScript := "install-run.ps1"
	if err := os.WriteFile(tempScript, []byte(psScript), 0644); err != nil {
		log.Fatalf("❌ Could not write PowerShell script: %v", err)
	}
	log.Printf("📝 Wrote script: %s", tempScript)

	// --- Run PowerShell script ---
	log.Println("🚀 Installing with PowerShell...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ PowerShell execution failed: %v", err)
	}

	log.Println("✅ Installation completed.")
}

// --- Maps user-friendly label to PowerShell function ---
func toInstallFunctionName(label string) string {
	label = strings.ToLower(label)
	label = strings.ReplaceAll(label, " ", "")
	label = strings.ReplaceAll(label, "-", "")

	switch label {
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
	default:
		return "Install-" + strings.Title(label)
	}
}

// --- Helpers for case-insensitive map access ---

func getCaseInsensitiveMap(data map[string]interface{}, target string) map[string]interface{} {
	for k, v := range data {
		if strings.EqualFold(k, target) {
			if subMap, ok := v.(map[string]interface{}); ok {
				return subMap
			}
		}
	}
	return nil
}

func getCaseInsensitiveList(data map[string]interface{}, target string) []string {
	for k, v := range data {
		if strings.EqualFold(k, target) {
			rawList, ok := v.([]interface{})
			if !ok {
				return nil
			}
			result := make([]string, 0, len(rawList))
			for _, item := range rawList {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

func getCaseInsensitiveString(data map[string]interface{}, target string) string {
	for k, v := range data {
		if strings.EqualFold(k, target) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}
