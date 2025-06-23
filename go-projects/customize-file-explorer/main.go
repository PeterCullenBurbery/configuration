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
	YamlKeyPath string            // e.g. "explorer.dark_mode"
	ValueMap    map[string]string // maps values like "on" or "true" to PowerShell functions
}

var mappings = []SettingMapping{
	// Explorer settings
	{"explorer.dark_mode", map[string]string{"true": "Set-DarkMode", "false": "Set-LightMode"}},
	{"explorer.search_box", map[string]string{"hidden": "Set-HideSearchBox", "shown": "Set-ShowSearchBox"}},
	{"explorer.file_extensions", map[string]string{"hidden": "Set-HideFileExtensions", "shown": "Set-ShowFileExtensions"}},
	{"explorer.hidden_files", map[string]string{"hidden": "Set-HideHiddenFiles", "shown": "Set-ShowHiddenFiles"}},
	{"explorer.start_menu_alignment", map[string]string{"left": "Set-StartMenuToLeft", "center": "Set-StartMenuToCenter"}},

	// Date/time settings
	{"date time settings.show seconds in taskbar", map[string]string{"on": "Set-ShowSecondsInTaskbar", "off": "Set-HideSecondsInTaskbar"}},
	{"date time settings.custom short date pattern", map[string]string{"on": "Set-CustomShortDatePattern", "off": "Reset-ShortDatePattern"}},
	{"date time settings.custom long date pattern", map[string]string{"on": "Set-CustomLongDatePattern", "off": "Reset-LongDatePattern"}},
	{"date time settings.custom time pattern", map[string]string{"on": "Set-CustomTimePattern", "off": "Reset-TimePatternToDefault"}},
	{"date time settings.24 hour time format", map[string]string{"on": "Set-24HourTimeFormat", "off": "Reset-12HourTimeFormat"}},
	{"date time settings.set first day of the week to monday", map[string]string{"on": "Set-FirstDayOfWeekMonday", "off": "Set-FirstDayOfWeekSunday"}},
}

func getNestedValue(m map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	var current interface{} = m
	for _, part := range parts {
		if mTyped, ok := current.(map[string]interface{}); ok {
			if val, exists := mTyped[part]; exists {
				current = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return current, true
}

func main() {
	configPath := flag.String("config", "", "Path to the config.yaml file (required)")
	modulePath := flag.String("module", "", "Path to the PowerShell module (.psm1) file (required)")
	logPath := flag.String("log", "", "Path to the log file (required)")
	flag.Parse()

	if *configPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ Error: --config, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

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

	var raw map[string]interface{}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	root, ok := raw["configuration_profile"].(map[string]interface{})
	if !ok {
		log.Fatalf("❌ 'configuration_profile' not found or invalid")
	}

	var psFunctions []string

	log.Println("🔧 Translating configuration to PowerShell functions...")
	for _, mapping := range mappings {
		if val, exists := getNestedValue(root, mapping.YamlKeyPath); exists {
			strVal := fmt.Sprintf("%v", val)
			if psFunc, ok := mapping.ValueMap[strings.ToLower(strVal)]; ok {
				psFunctions = append(psFunctions, psFunc)
				log.Printf("✔️  %s = %s → %s", mapping.YamlKeyPath, strVal, psFunc)
			} else {
				log.Printf("⚠️  Unknown value: %s = %s", mapping.YamlKeyPath, strVal)
			}
		} else {
			log.Printf("⚠️  Key not found in YAML: %s", mapping.YamlKeyPath)
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
	defer os.Remove(tempScript)

	log.Printf("📝 Wrote script: %s", tempScript)

	log.Println("🚀 Running PowerShell script...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Run()
	if err != nil {
		log.Fatalf("❌ PowerShell execution failed: %v", err)
	}
	log.Println("✅ PowerShell script executed successfully.")
}
