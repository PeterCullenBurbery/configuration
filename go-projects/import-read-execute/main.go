package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
	// Allow optional path to config.yaml as CLI arg
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Load and parse YAML file
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read %s: %v", configPath, err)
	}

	var raw map[string]map[string]interface{}
	err = yaml.Unmarshal(content, &raw)
	if err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	config := raw["configuration_profile"]
	var psFunctions []string

	for _, mapping := range mappings {
		if val, exists := config[mapping.YamlKey]; exists {
			strVal := fmt.Sprintf("%v", val)
			if psFunc, ok := mapping.ValueMap[strVal]; ok {
				psFunctions = append(psFunctions, psFunc)
			} else {
				log.Printf("⚠️ No match for key %q with value %q", mapping.YamlKey, strVal)
			}
		}
	}

	if len(psFunctions) == 0 {
		log.Println("⚠️ No functions to execute from config.")
		return
	}

	psScript := "Import-Module ./output/MyModule.psm1\n"
	for _, fn := range psFunctions {
		psScript += fn + "\n"
	}

	tempScript := "run.ps1"
	err = os.WriteFile(tempScript, []byte(psScript), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write PowerShell script: %v", err)
	}

	fmt.Printf("🚀 Running PowerShell with config %s...\n", configPath)
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ PowerShell execution failed: %v", err)
	}
	fmt.Println("✅ Configuration complete.")
}
