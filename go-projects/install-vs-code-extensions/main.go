package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type ExtensionsConfig struct {
	VSCodeExtensions []string `yaml:"vs_code_extensions"`
}

func main() {
	yamlPath := flag.String("yaml", "", "Path to vs-code-extensions.yaml (required)")
	flag.Parse()

	if *yamlPath == "" {
		log.Fatal("❌ --yaml is required")
	}

	data, err := os.ReadFile(*yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML file: %v", err)
	}

	var config ExtensionsConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	if len(config.VSCodeExtensions) == 0 {
		log.Println("⚠️ No extensions found in vs_code_extensions list.")
		return
	}

	for _, ext := range config.VSCodeExtensions {
		log.Printf("🔧 Installing VS Code extension: %s\n", ext)

		cmd := exec.Command("code", "--install-extension", ext)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Printf("❌ Failed to install %s: %v\n", ext, err)
		} else {
			log.Printf("✅ Installed: %s\n", ext)
		}
	}
}