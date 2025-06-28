package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PeterCullenBurbery/go_functions_002/system_management_functions"
)

type ProgramEntry struct {
	Name         string   `yaml:"name"`
	Alternatives []string `yaml:"alternatives"`
	WingetID     string   `yaml:"winget id,omitempty"`
	ChocoID      string   `yaml:"choco id,omitempty"`
}

type InstallYaml struct {
	Install map[string]map[string]ProgramEntry `yaml:"install"`
}

func main() {
	whatPath := flag.String("what", "", "Path to what-to-install.yaml (required)")
	installPath := flag.String("install", "", "Path to install.yaml (required)")
	logPath := flag.String("log", "", "Path to log file (required)")
	flag.Parse()

	if *whatPath == "" || *installPath == "" || *logPath == "" {
		fmt.Println("❌ --what, --install, and --log are required.")
		flag.Usage()
		os.Exit(1)
	}

	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Load install.yaml
	var installData InstallYaml
	rawInstallData, err := os.ReadFile(*installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read install.yaml: %v", err)
	}
	if err := yaml.Unmarshal(rawInstallData, &installData); err != nil {
		log.Fatalf("❌ Failed to parse install.yaml: %v", err)
	}

	// Build lookup maps
	altToCanonical := make(map[string]string)
	canonicalToMeta := make(map[string]ProgramEntry)
	canonicalToCategory := make(map[string]string)

	for category, programs := range installData.Install {
		for canonical, meta := range programs {
			canonicalTrimmed := strings.TrimSpace(canonical)
			canonicalToMeta[canonicalTrimmed] = meta
			canonicalToCategory[canonicalTrimmed] = category
			altToCanonical[strings.ToLower(canonicalTrimmed)] = canonicalTrimmed
			for _, alt := range meta.Alternatives {
				altToCanonical[strings.ToLower(strings.TrimSpace(alt))] = canonicalTrimmed
			}
		}
	}

	// Load what-to-install.yaml
	whatData := make(map[string]interface{})
	rawWhatData, err := os.ReadFile(*whatPath)
	if err != nil {
		log.Fatalf("❌ Failed to read what-to-install.yaml: %v", err)
	}
	if err := yaml.Unmarshal(rawWhatData, &whatData); err != nil {
		log.Fatalf("❌ Failed to parse what-to-install.yaml: %v", err)
	}

	installSection := getCaseInsensitiveMap(whatData, "install")
	if installSection == nil {
		log.Fatal("❌ Missing 'install' section in what-to-install.yaml.")
	}
	requested := getCaseInsensitiveList(installSection, "programs to install")

	// Process programs
	for _, req := range requested {
		lookup := strings.ToLower(strings.TrimSpace(req))
		canonical, ok := altToCanonical[lookup]
		if !ok {
			log.Printf("❌ Unsupported program: %s (skipped)", req)
			continue
		}

		category := canonicalToCategory[canonical]
		meta := canonicalToMeta[canonical]

		log.Printf("✅ Supported program: %s → %s (category: %s)", req, canonical, category)

		switch category {
		case "automatically installed":
			log.Printf("ℹ️  %s is already installed automatically. Skipping.", canonical)

		case "winget":
			if meta.WingetID == "" {
				log.Printf("⚠️ Missing Winget ID for %s", canonical)
				continue
			}
			err := system_management_functions.Winget_install(canonical, meta.WingetID)
			if err != nil {
				log.Printf("❌ Winget install failed for %s: %v", canonical, err)
			} else {
				log.Printf("✅ Installed %s via Winget.", canonical)
			}

		case "choco":
			if meta.ChocoID == "" {
				log.Printf("⚠️ Missing Choco ID for %s", canonical)
				continue
			}
			err := system_management_functions.Choco_install(meta.ChocoID)
			if err != nil {
				log.Printf("❌ Chocolatey install failed for %s: %v", canonical, err)
			} else {
				log.Printf("✅ Installed %s via Chocolatey.", canonical)
			}

		default:
			log.Printf("⚠️ Unknown or unhandled category '%s' for %s", category, canonical)
		}
	}

	log.Println("🎉 Installation process finished.")
}