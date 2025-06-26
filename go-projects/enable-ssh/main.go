package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigurationProfile struct {
		SSH string `yaml:"ssh"`
	} `yaml:"configuration_profile"`
}

func main() {
	yamlPath := flag.String("yaml", "", "Path to config.yaml (required)")
	modulePath := flag.String("module", "", "Path to PowerShell module (.psm1) (required)")
	logPath := flag.String("log", "", "Path to log file (required)")
	flag.Parse()

	if *yamlPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ --yaml, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

	// Open log file
	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Load YAML
	data, err := os.ReadFile(*yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("❌ Failed to parse YAML file: %v", err)
	}

	if strings.ToLower(strings.TrimSpace(cfg.ConfigurationProfile.SSH)) != "on" {
		log.Println("ℹ️ SSH is disabled in config.yaml. Skipping SSH setup.")
		return
	}

	// Generate PowerShell script content
	script := fmt.Sprintf(`Import-Module '%s'
Enable-SSH
Enable-SSHFirewallRule
`, *modulePath)

	scriptName := "enable-ssh.ps1"
	if err := os.WriteFile(scriptName, []byte(script), 0644); err != nil {
		log.Fatalf("❌ Failed to write PowerShell script: %v", err)
	}

	// Execute PowerShell script with -NoProfile to avoid profile interference
	log.Println("🔐 Running Enable-SSH and Enable-SSHFirewallRule...")
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptName)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	start := time.Now()
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ SSH setup failed: %v", err)
	}
	duration := time.Since(start).Seconds()
	log.Printf("✅ SSH setup completed in %.2f seconds.", duration)
}