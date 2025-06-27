package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/PeterCullenBurbery/go-functions"
)

func main() {
	yamlPath := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\install.yaml`

	// Read YAML file
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	// Get downloads section
	install := getCaseInsensitiveMap(raw, "install")
	if install == nil {
		log.Fatal("❌ No 'install' section found.")
	}

	downloads := getCaseInsensitiveMap(install, "downloads")
	if downloads == nil {
		log.Fatal("❌ No 'downloads' section found.")
	}

	globalDir := strings.TrimSpace(getNestedString(downloads, "global download directory"))
	perApp := getNestedMap(downloads, "per app download directories")

	subDownload := strings.TrimSpace(getCaseInsensitiveString(perApp, "nirsoft"))
	if globalDir == "" || subDownload == "" {
		log.Fatal("❌ Could not resolve nirsoft download path.")
	}

	// Generate timestamp
	rawTimestamp, err := gofunctions.DateTimeStamp()
	if err != nil {
		log.Fatalf("❌ Failed to generate timestamp: %v", err)
	}
	safeTimestamp := gofunctions.SafeTimeStamp(rawTimestamp, 1)

	// Final path
	nirsoftPath := filepath.Join(globalDir, subDownload, safeTimestamp)
	log.Printf("📁 Creating directory: %s", nirsoftPath)

	if err := os.MkdirAll(nirsoftPath, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create directory: %v", err)
	}

	log.Println("✅ Nirsoft timestamped folder created.")
}

func getCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if submap, ok := v.(map[string]interface{}); ok {
				return submap
			}
		}
	}
	return nil
}

func getCaseInsensitiveString(m map[string]interface{}, key string) string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}

func getNestedString(m map[string]interface{}, key string) string {
	if val := getCaseInsensitiveString(m, key); val != "" {
		return val
	}
	return ""
}

func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	if val := getCaseInsensitiveMap(m, key); val != nil {
		return val
	}
	return nil
}
