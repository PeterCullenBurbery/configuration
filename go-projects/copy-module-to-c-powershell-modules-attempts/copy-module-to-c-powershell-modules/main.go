package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// getCaseInsensitiveMap returns a nested map for a key, ignoring case
func getCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if subMap, ok := v.(map[string]interface{}); ok {
				return subMap
			}
		}
	}
	return nil
}

// getCaseInsensitiveString returns a string value from a map, ignoring case
func getCaseInsensitiveString(m map[string]interface{}, key string) string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if val, ok := v.(string); ok {
				return strings.TrimSpace(val)
			}
		}
	}
	return ""
}

func main() {
	// --- Command line flags ---
	psm1Path := flag.String("psm1", "", "Path to .psm1 file (required)")
	psd1Path := flag.String("psd1", "", "Path to .psd1 file (required)")
	yamlPath := flag.String("yaml", "", "Path to YAML config (required)")
	moduleName := flag.String("name", "", "Module name (required)")
	flag.Parse()

	// --- Validate input ---
	if *psm1Path == "" || *psd1Path == "" || *yamlPath == "" || *moduleName == "" {
		flag.Usage()
		log.Fatal("❌ All arguments --psm1, --psd1, --yaml, and --name are required.")
	}

	// --- Read and parse YAML ---
	yamlData, err := os.ReadFile(*yamlPath)
	if err != nil {
		log.Fatalf("❌ Could not read YAML file: %v", err)
	}

	var root map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &root); err != nil {
		log.Fatalf("❌ YAML parsing error: %v", err)
	}

	configProfile := getCaseInsensitiveMap(root, "configuration_profile")
	if configProfile == nil {
		log.Fatal("❌ Missing 'configuration_profile' section in YAML.")
	}

	psModulePath := getCaseInsensitiveString(configProfile, "powershell modules")
	if psModulePath == "" {
		log.Fatal("❌ Could not extract 'powershell modules' from YAML.")
	}

	// --- Inline PowerShell function ---
	psFunction := `
function Install-CustomModule {
    [CmdletBinding()]
    param (
        [Parameter(Position = 0, Mandatory = $true)][string]$psm1_path,
        [Parameter(Position = 1, Mandatory = $true)][string]$psd1_path,
        [Parameter(Position = 2, Mandatory = $true)][string]$module_name,
        [Parameter(Position = 3, Mandatory = $true)][string]$target_directory
    )
    $module_folder = Join-Path -Path $target_directory -ChildPath $module_name
    if (-not (Test-Path -Path $module_folder)) {
        New-Item -Path $module_folder -ItemType Directory -Force | Out-Null
        Write-Host "📁 Created folder: $module_folder"
    }
    try {
        Copy-Item -Path $psm1_path -Destination (Join-Path $module_folder "$module_name.psm1") -Force
        Copy-Item -Path $psd1_path -Destination (Join-Path $module_folder "$module_name.psd1") -Force
        Write-Host "✅ Module '$module_name' installed to: $module_folder"
    } catch {
        Write-Error "❌ Failed to install module: $_"
        return $false
    }
    return $true
}`

	// --- Complete PowerShell command ---
	psCommand := fmt.Sprintf(`%s; Install-CustomModule -psm1_path "%s" -psd1_path "%s" -module_name "%s" -target_directory "%s"`,
		psFunction, *psm1Path, *psd1Path, *moduleName, psModulePath)

	fmt.Println("🚀 Executing PowerShell command:")
	fmt.Println(psCommand)

	cmd := exec.Command("pwsh", "-NoProfile", "-Command", psCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ PowerShell execution failed: %v", err)
	}

	fmt.Println("✅ Module installation completed.")
}