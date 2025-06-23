package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	installPath := flag.String("install", "", "Path to install.yaml (required)")
	modulePath := flag.String("module", "", "Path to PowerShell module (.psm1) (required)")
	logPath := flag.String("log", "", "Path to runtime execution log file (required)")
	flag.Parse()

	if *installPath == "" || *modulePath == "" || *logPath == "" {
		fmt.Println("❌ --install, --module, and --log are all required.")
		flag.Usage()
		os.Exit(1)
	}

	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	var raw map[string]interface{}
	data, err := os.ReadFile(*installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	installSection := getCaseInsensitiveMap(raw, "install")
	if installSection == nil {
		log.Fatal("❌ Missing 'install' section.")
	}

	programs := getCaseInsensitiveList(installSection, "programs to install")
	logs := getCaseInsensitiveMap(installSection, "logs")
	downloads := getCaseInsensitiveMap(installSection, "downloads")

	globalLogDir := strings.TrimSpace(getNestedString(logs, "global log directory"))
	perAppLogs := getNestedMap(logs, "per app log directories")
	globalDownloadDir := strings.TrimSpace(getNestedString(downloads, "global download directory"))
	perAppDownloads := getNestedMap(downloads, "per app download directories")

	if globalLogDir == "" || globalDownloadDir == "" {
		log.Fatal("❌ Missing 'global log directory' or 'global download directory'.")
	}

	_ = os.MkdirAll(globalLogDir, os.ModePerm)
	_ = os.MkdirAll(globalDownloadDir, os.ModePerm)

	var psScript strings.Builder
	psScript.WriteString(fmt.Sprintf("Import-Module '%s'\n", *modulePath))

	for _, label := range programs {
		funcName := toInstallFunctionName(label)
		log.Printf("✔️ Queued installer: %s → %s", label, funcName)

		if strings.EqualFold(funcName, "Install-CherryTree") {
			appKey := "cherry tree"
			subLog := strings.TrimSpace(getCaseInsensitiveString(perAppLogs, appKey))
			subDownload := strings.TrimSpace(getCaseInsensitiveString(perAppDownloads, appKey))

			timestamp := formatTimestamp()
			logDir := filepath.Join(globalLogDir, subLog)
			logFileName := fmt.Sprintf("cherrytree_%s.log", timestamp)
			cherryLogPath := filepath.Join(logDir, logFileName)
			cherryInstallPath := filepath.Join(globalDownloadDir, subDownload)
			installerPath := filepath.Join(cherryInstallPath, "cherrytree_1.5.0.0_win64_setup.exe")
			installerURL := "https://www.giuspen.net/software/cherrytree_1.5.0.0_win64_setup.exe"

			_ = os.MkdirAll(logDir, os.ModePerm)
			_ = os.MkdirAll(cherryInstallPath, os.ModePerm)

			if !fileExists(installerPath) {
				log.Printf("🌐 Downloading CherryTree from: %s", installerURL)
				if err := downloadFile(installerPath, installerURL); err != nil {
					log.Fatalf("❌ Download failed: %v", err)
				}
				log.Println("✅ Downloaded CherryTree.")
			} else {
				log.Println("📁 CherryTree installer already present.")
			}

			log.Printf("📝 CherryTree log path: %s", cherryLogPath)
			psScript.WriteString(fmt.Sprintf(`%s -log '%s' -installPath '%s'`+"\n", funcName, cherryLogPath, cherryInstallPath))
		} else {
			psScript.WriteString(funcName + "\n")
		}
	}

	tempScript := "install-run.ps1"
	if err := os.WriteFile(tempScript, []byte(psScript.String()), 0644); err != nil {
		log.Fatalf("❌ Failed to write PowerShell script: %v", err)
	}

	log.Println("🚀 Executing install script...")
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScript)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ PowerShell script failed: %v", err)
	}
	log.Println("✅ Installation complete.")
}

// --- Helpers ---

func toInstallFunctionName(label string) string {
	l := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(label, " ", ""), "-", ""))
	switch l {
	case "powershell7":
		return "Install-PowerShell-7"
	case "vscode":
		return "Install-VSCode"
	case "7zip":
		return "Install-7Zip"
	case "voidtoolseverything":
		return "Install-Voidtools-Everything"
	case "winscp":
		return "Install-WinSCP"
	case "mobaxterm":
		return "Install-MobaXterm"
	case "choco", "chocolatey":
		return "Install-Choco"
	case "cherrytree":
		return "Install-CherryTree"
	case "go":
		return "Install-Go"
	default:
		return "Install-" + strings.Title(l)
	}
}

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
