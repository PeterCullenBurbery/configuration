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

type SettingMapping struct {
	YamlKey  string
	ValueMap map[string]string
}

var mappings = []SettingMapping{
	{
		YamlKey: "dark_mode",
		ValueMap: map[string]string{
			"true":  "Set-DarkMode",
			"false": "Set-LightMode",
		},
	},
	{
		YamlKey: "search_box",
		ValueMap: map[string]string{
			"hidden": "Set-HideSearchBox",
			"shown":  "Set-ShowSearchBox",
		},
	},
	{
		YamlKey: "file_extensions",
		ValueMap: map[string]string{
			"hidden": "Set-HideFileExtensions",
			"shown":  "Set-ShowFileExtensions",
		},
	},
	{
		YamlKey: "hidden_files",
		ValueMap: map[string]string{
			"hidden": "Set-HideHiddenFiles",
			"shown":  "Set-ShowHiddenFiles",
		},
	},
	{
		YamlKey: "start_menu_alignment",
		ValueMap: map[string]string{
			"left":   "Set-StartMenuToLeft",
			"center": "Set-StartMenuToCenter",
		},
	},
}

func main() {
	// Define required flags
	configPath := flag.String("config", "", "Path to the config.yaml file (required)")
	modulePath := flag.String("module", "", "Path to the PowerShell module (.psm1) file (required)")
	logPath := flag.String("log", "", "Path to the log file (required)")
	flag.Parse()

	// Enforce required arguments
	if *configPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ Error: --config, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

	// Open log file and set output
	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file %s: %v\n", *logPath, err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	log.Println("📄 Reading config:", *configPath)
	content, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	var raw map[string]map[string]interface{}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	config := raw["configuration_profile"]
	var psFunctions []string

	log.Println("🔧 Translating configuration to PowerShell functions...")
	for _, mapping := range mappings {
		if val, exists := config[mapping.YamlKey]; exists {
			strVal := fmt.Sprintf("%v", val)
			if psFunc, ok := mapping.ValueMap[strings.ToLower(strVal)]; ok {
				psFunctions = append(psFunctions, psFunc)
				log.Printf("✔️  %s = %s → %s", mapping.YamlKey, strVal, psFunc)
			} else {
				log.Printf("⚠️  Unknown value: %s = %s", mapping.YamlKey, strVal)
			}
		}
	}

	if len(psFunctions) == 0 {
		log.Println("⚠️ No PowerShell functions matched. Exiting.")
		return
	}

	psScript := fmt.Sprintf("Import-Module '%s'\n", *modulePath)
	for _, fn := range psFunctions {
		psScript += fn + "\n"
	}

	tempScript := "run.ps1"
	if err := os.WriteFile(tempScript, []byte(psScript), 0644); err != nil {
		log.Fatalf("❌ Failed to write temporary PowerShell script: %v", err)
	}
	log.Printf("📝 Wrote script: %s", tempScript)

	log.Println("🚀 Running PowerShell script...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = logFile // capture only in log
	cmd.Stderr = logFile
	err = cmd.Run()
	if err != nil {
		log.Fatalf("❌ PowerShell execution failed: %v", err)
	}
	log.Println("✅ PowerShell script executed successfully.")
}
