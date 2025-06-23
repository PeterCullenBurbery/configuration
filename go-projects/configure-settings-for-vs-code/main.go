// edit_vscode_settings.go

// To see VS Code settings, use:
// PowerShell: code $env:APPDATA\Code\User\settings.json

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Get VS Code settings.json path
	path, err := settingsPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Read existing settings.json if present
	var cfg map[string]interface{}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cfg = make(map[string]interface{}) // No file, start fresh
		} else {
			fmt.Fprintf(os.Stderr, "failed to read settings.json: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := json.Unmarshal(data, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse existing JSON: %v\n", err)
			os.Exit(1)
		}
	}

	// Backup original if it existed
	if data != nil {
		if err := backup(path, data); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to backup original settings.json: %v\n", err)
		}
	}

	// Get full Desktop path (e.g., C:\Users\peter\Desktop)
	desktopPath, err := getResolvedDesktopPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine desktop path: %v\n", err)
		os.Exit(1)
	}

	// Desired settings to insert/overwrite
	cfg["files.autoSave"] = "afterDelay"
	cfg["powershell.cwd"] = desktopPath
	cfg["terminal.integrated.cwd"] = desktopPath
	cfg["terminal.integrated.enableMultiLinePasteWarning"] = "never"
	cfg["terminal.integrated.persistentSessionScrollback"] = 10000000
	cfg["terminal.integrated.rightClickBehavior"] = "default"
	cfg["terminal.integrated.scrollback"] = 10000000
	cfg["workbench.startupEditor"] = "none"
	cfg["explorer.confirmDragAndDrop"] = false
	cfg["explorer.confirmDelete"] = false
	cfg["redhat.telemetry.enabled"] = true
	cfg["editor.renderWhitespace"] = "all"

	// YAML-specific editor settings
	cfg["[yaml]"] = map[string]interface{}{
		"editor.insertSpaces":      true,
		"editor.tabSize":           2,
		"editor.detectIndentation": false,
	}

	// Marshal JSON with indentation
	out, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal merged JSON: %v\n", err)
		os.Exit(1)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create directory %q: %v\n", dir, err)
		os.Exit(1)
	}

	// Atomic write: temp file -> rename
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write temp settings.json: %v\n", err)
		os.Exit(1)
	}
	if err := os.Rename(tmp, path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to overwrite settings.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ VS Code settings.json updated successfully.")
}

func settingsPath() (string, error) {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return "", errors.New("APPDATA environment variable not set")
	}
	return filepath.Join(appdata, "Code", "User", "settings.json"), nil
}

func backup(origPath string, data []byte) error {
	dir := filepath.Dir(origPath)
	base := filepath.Base(origPath)
	ts := time.Now().Format("20060102_150405")
	bakName := fmt.Sprintf("%s.bak.%s", base, ts)
	bakPath := filepath.Join(dir, bakName)
	return os.WriteFile(bakPath, data, 0o644)
}

func getResolvedDesktopPath() (string, error) {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return "", errors.New("USERPROFILE environment variable not set")
	}
	desktop := filepath.Join(userProfile, "Desktop")
	return desktop, nil
}
