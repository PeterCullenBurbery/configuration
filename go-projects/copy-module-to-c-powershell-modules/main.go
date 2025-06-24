package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigurationProfile struct {
		PowerShellModules string `yaml:"powershell modules"`
	} `yaml:"configuration_profile"`
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func fileMustExist(path string, label string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("❌ %s file does not exist: %s", label, path)
	}
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: <program> <psm1 path> <psd1 path> <yaml path> <module name>")
		os.Exit(1)
	}

	psm1Path := os.Args[1]
	psd1Path := os.Args[2]
	yamlPath := os.Args[3]
	moduleName := os.Args[4]

	// --- Validate file existence ---
	fileMustExist(psm1Path, ".psm1")
	fileMustExist(psd1Path, ".psd1")
	fileMustExist(yamlPath, "YAML")

	// --- Read and parse YAML ---
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	// --- Clean and normalize base module directory ---
	baseModuleDir := filepath.Clean(strings.TrimSpace(config.ConfigurationProfile.PowerShellModules))
	fullModuleDir := filepath.Join(baseModuleDir, moduleName)

	// --- Create base and module directories ---
	if err := os.MkdirAll(fullModuleDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create module directory: %v", err)
	}
	fmt.Printf("📁 Created module directory: %s\n", fullModuleDir)

	// --- Copy .psm1 ---
	psm1Dst := filepath.Join(fullModuleDir, filepath.Base(psm1Path))
	if err := copyFile(psm1Path, psm1Dst); err != nil {
		log.Fatalf("❌ Failed to copy .psm1: %v", err)
	}
	fmt.Printf("✅ Copied: %s → %s\n", psm1Path, psm1Dst)

	// --- Copy .psd1 ---
	psd1Dst := filepath.Join(fullModuleDir, filepath.Base(psd1Path))
	if err := copyFile(psd1Path, psd1Dst); err != nil {
		log.Fatalf("❌ Failed to copy .psd1: %v", err)
	}
	fmt.Printf("✅ Copied: %s → %s\n", psd1Path, psd1Dst)
}