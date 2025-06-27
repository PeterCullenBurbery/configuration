package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// --- Helper functions ---

func getCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if result, ok := v.(map[string]interface{}); ok {
				return result
			}
		}
	}
	return nil
}

func getCaseInsensitiveList(m map[string]interface{}, key string) []string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			raw, ok := v.([]interface{})
			if !ok {
				return nil
			}
			var result []string
			for _, val := range raw {
				if s, ok := val.(string); ok {
					result = append(result, s)
				}
			}
			return result
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
	if sub := getCaseInsensitiveMap(m, key); sub != nil {
		for _, v := range sub {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	return getCaseInsensitiveMap(m, key)
}

func downloadFile(dest, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func formatTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// New helper to run each individual script:
func runPowerShellScript(filename, content string, logFile *os.File) {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("❌ Failed to write script %s: %v", filename, err)
	}
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", filename)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to run %s: %v", filename, err)
	}
	log.Printf("✅ Finished: %s", filename)
}
