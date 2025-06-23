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
	ValueMap map[string]string // map from YAML value to PowerShell function
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
	// Step 1: Read the config.yaml
	configFile := "config.yaml"
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("❌ Failed to read %s: %v", configFile, err)
	}

	// Step 2: Unmarshal into a generic map
	var raw map[string]map[string]interface{}
	err = yaml.Unmarshal(content, &raw)
	if err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	config := raw["configuration_profile"]
	var psFunctions []string

	// Step 3: Process config using mappings
	for _, mapping := range mappings {
		if val, exists := config[mapping.YamlKey]; exists {
			strVal := fmt.Sprintf("%v", val)
			if psFunc, ok := mapping.ValueMap[strVal]; ok {
				psFunctions = append(psFunctions, psFunc)
			} else {
				log.Printf("⚠️  No PowerShell mapping for key=%q value=%q", mapping.YamlKey, strVal)
			}
		}
	}

	if len(psFunctions) == 0 {
		log.Println("⚠️  No matching PowerShell functions found. Nothing to do.")
		return
	}

	// Step 4: Create PowerShell script
	psScript := "Import-Module ./output/MyModule.psm1\n"
	for _, fn := range psFunctions {
		psScript += fn + "\n"
	}

	// Step 5: Write the script to file
	tempScript := "run.ps1"
	err = os.WriteFile(tempScript, []byte(psScript), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write %s: %v", tempScript, err)
	}

	// Step 6: Run the script
	fmt.Println("🚀 Running configuration script...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ PowerShell script execution failed: %v", err)
	}
	fmt.Println("✅ Configuration applied successfully.")
}
