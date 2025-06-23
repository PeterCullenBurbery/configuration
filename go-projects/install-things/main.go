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

type InstallFile struct {
	Install map[string]interface{} `yaml:"install"`
}

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

	var installConfig InstallFile
	if err := yaml.Unmarshal(data, &installConfig); err != nil {
		log.Fatalf("❌ Failed to parse install.yaml: %v", err)
	}

	// --- Map keys to PowerShell Install-* functions ---
	var psFunctions []string
	for key := range installConfig.Install {
		funcName := toInstallFunctionName(key)
		psFunctions = append(psFunctions, funcName)
		log.Printf("✔️  Queued installer: %s → %s", key, funcName)
	}

	if len(psFunctions) == 0 {
		log.Println("⚠️ No functions listed in install.yaml. Exiting.")
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

// Converts "VS code" → "Install-VSCode", "7zip" → "Install-7Zip"
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
	default:
		return "Install-" + strings.Title(label)
	}
}
