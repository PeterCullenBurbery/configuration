package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigurationProfile struct {
		PowerShellModules string `yaml:"powershell modules"`
	} `yaml:"configuration_profile"`
}

func readPowerShellModulePath(yamlPath string) (string, error) {
	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return "", err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	modulePath := strings.TrimSpace(cfg.ConfigurationProfile.PowerShellModules)
	if modulePath == "" {
		return "", fmt.Errorf("no powershell modules path found in YAML")
	}
	return modulePath, nil
}

func runGoCLI(goCLI, modulePath string) error {
	cmd := exec.Command(goCLI, "add-topsmodulepath", modulePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Define flags
	cliFlag := flag.String("cli", "", "Path to Go CLI executable")
	yamlFlag := flag.String("yaml", "", "Path to YAML configuration file")
	flag.Parse()

	// Fallback to positionals if flags not used
	var cliPath, yamlPath string

	if *cliFlag != "" && *yamlFlag != "" {
		cliPath = *cliFlag
		yamlPath = *yamlFlag
	} else if flag.NArg() == 2 {
		cliPath = flag.Arg(0)
		yamlPath = flag.Arg(1)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  1234.exe --cli <go-cli> --yaml <config.yaml>")
		fmt.Println("  1234.exe <go-cli> <config.yaml>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	modulePath, err := readPowerShellModulePath(yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v\n", err)
	}

	fmt.Printf("🚀 Running: %s add-topsmodulepath %s\n", cliPath, modulePath)

	if err := runGoCLI(cliPath, modulePath); err != nil {
		log.Fatalf("❌ Failed to run Go CLI: %v\n", err)
	}
}