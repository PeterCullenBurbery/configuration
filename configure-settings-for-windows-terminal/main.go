// edit_windows_terminal.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GUID constants
const (
	ps7GUID   = "{574e775e-4f2a-5b96-ac1e-a2962a402336}"
	winPSGUID = "{61c54bbd-c2c6-5271-96e7-009a87ff44bf}"
	cmdGUID   = "{0caa0dad-35be-5f56-a8ff-afceeeaa6101}"
	azureGUID = "{b453ae62-4e3d-5e58-b989-0a998ec441b8}"
)

func main() {
	path, err := settingsPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// Read original
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read settings.json: %v\n", err)
		os.Exit(1)
	}
	// Backup
	if err := backup(path, data); err != nil {
		fmt.Fprintf(os.Stderr, "warning: backup failed: %v\n", err)
		// continue anyway
	}
	// Unmarshal
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse JSON: %v\n", err)
		os.Exit(1)
	}
	// Transform
	if err := transform(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "transformation error: %v\n", err)
		os.Exit(1)
	}
	// Marshal with indentation
	out, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	// Write back
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write temp file: %v\n", err)
		os.Exit(1)
	}
	// Replace original
	if err := os.Rename(tmpPath, path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to overwrite settings.json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("settings.json updated successfully (backup created).")
}

// settingsPath builds the path to settings.json using LOCALAPPDATA.
func settingsPath() (string, error) {
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		return "", errors.New("LOCALAPPDATA not set")
	}
	// Windows Terminal package path
	rel := filepath.Join(
		"Packages",
		"Microsoft.WindowsTerminal_8wekyb3d8bbwe",
		"LocalState",
		"settings.json",
	)
	full := filepath.Join(local, rel)
	return full, nil
}

// backup writes the original data to settings.json.bak.TIMESTAMP
func backup(origPath string, data []byte) error {
	dir := filepath.Dir(origPath)
	base := filepath.Base(origPath)
	ts := time.Now().Format("20060102_150405")
	bakName := fmt.Sprintf("%s.bak.%s", base, ts)
	bakPath := filepath.Join(dir, bakName)
	return os.WriteFile(bakPath, data, 0o644)
}

// transform applies the JSON modifications in-place on cfg.
func transform(cfg map[string]interface{}) error {
	// 1. Set defaultProfile
	cfg["defaultProfile"] = ps7GUID

	// 2. Locate profiles object
	profilesRaw, ok := cfg["profiles"]
	if !ok {
		return errors.New(`missing "profiles" object`)
	}
	profiles, ok := profilesRaw.(map[string]interface{})
	if !ok {
		return errors.New(`"profiles" is not an object`)
	}
	// 3. Set profiles.defaults
	profiles["defaults"] = map[string]interface{}{
		"elevate":     true,
		"historySize": 1000000000,
	}
	// 4. Rebuild profiles.list
	listRaw, ok := profiles["list"]
	if !ok {
		return errors.New(`missing "profiles.list"`)
	}
	listSlice, ok := listRaw.([]interface{})
	if !ok {
		return errors.New(`"profiles.list" is not an array`)
	}
	// Index existing entries by GUID
	entries := make(map[string]map[string]interface{})
	for _, item := range listSlice {
		if m, ok := item.(map[string]interface{}); ok {
			if g, ok := m["guid"].(string); ok {
				entries[g] = m
			}
		}
	}
	// Build new list in order: PS7, Windows PowerShell, Command Prompt, Azure Cloud Shell
	var newList []interface{}

	// helper to prepare each entry
	prep := func(guid string) map[string]interface{} {
		if e, exists := entries[guid]; exists {
			return e
		}
		return make(map[string]interface{})
	}

	// 4a. PowerShell 7
	ePS7 := prep(ps7GUID)
	ePS7["guid"] = ps7GUID
	ePS7["name"] = "PowerShell 7"
	ePS7["hidden"] = false
	ePS7["source"] = "Windows.Terminal.PowershellCore"
	delete(ePS7, "commandline")
	newList = append(newList, ePS7)

	// 4b. Windows PowerShell -> rename to PowerShell 5
	eWin := prep(winPSGUID)
	eWin["guid"] = winPSGUID
	eWin["name"] = "PowerShell 5"
	eWin["hidden"] = false
	// ensure commandline exists; if absent, set default
	if _, ok := eWin["commandline"]; !ok {
		eWin["commandline"] = "%SystemRoot%\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	}
	delete(eWin, "source")
	newList = append(newList, eWin)

	// 4c. Command Prompt
	eCmd := prep(cmdGUID)
	eCmd["guid"] = cmdGUID
	eCmd["name"] = "Command Prompt"
	eCmd["hidden"] = false
	if _, ok := eCmd["commandline"]; !ok {
		eCmd["commandline"] = "%SystemRoot%\\System32\\cmd.exe"
	}
	delete(eCmd, "source")
	newList = append(newList, eCmd)

	// 4d. Azure Cloud Shell
	eAz := prep(azureGUID)
	eAz["guid"] = azureGUID
	eAz["name"] = "Azure Cloud Shell"
	eAz["hidden"] = false
	eAz["source"] = "Windows.Terminal.Azure"
	delete(eAz, "commandline")
	newList = append(newList, eAz)

	profiles["list"] = newList
	return nil
}
