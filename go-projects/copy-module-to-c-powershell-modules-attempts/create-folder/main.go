package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigurationProfile struct {
		PowerShellModules string `yaml:"powershell modules"`
	} `yaml:"configuration_profile"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: create_module_folder <path to config.yaml>")
		os.Exit(1)
	}

	configPath := os.Args[1]

	yamlData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	moduleDir := filepath.Clean(config.ConfigurationProfile.PowerShellModules)

	err = os.MkdirAll(moduleDir, 0755)
	if err != nil {
		log.Fatalf("❌ Failed to create folder '%s': %v", moduleDir, err)
	}

	fmt.Printf("✅ Created folder: %s\n", moduleDir)
}
