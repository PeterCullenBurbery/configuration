package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Keybinding struct {
	Key     string `json:"key"`
	Command string `json:"command"`
	When    string `json:"when,omitempty"`
}

func main() {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		fmt.Println("❌ APPDATA environment variable not set.")
		return
	}

	keybindingsPath := filepath.Join(appData, "Code", "User", "keybindings.json")

	var existingBindings []Keybinding

	// Check if file exists
	if _, err := os.Stat(keybindingsPath); err == nil {
		data, err := os.ReadFile(keybindingsPath)
		if err == nil && len(data) > 0 {
			if err := json.Unmarshal(data, &existingBindings); err != nil {
				fmt.Printf("❌ Failed to parse existing keybindings.json: %v\n", err)
				return
			}
		}
	} else if os.IsNotExist(err) {
		// File doesn't exist: ensure parent folder exists
		dir := filepath.Dir(keybindingsPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ Failed to create directory %s: %v\n", dir, err)
			return
		}
		// Continue with empty keybindings list
		fmt.Println("ℹ️ keybindings.json does not exist, creating a new one.")
	} else {
		fmt.Printf("❌ Failed to stat keybindings.json: %v\n", err)
		return
	}

	// Desired keybindings
	newBindings := []Keybinding{
		{
			Key:     "ctrl+a",
			Command: "workbench.action.terminal.selectAll",
			When:    "terminalFocus",
		},
		{
			Key:     "ctrl+shift+a",
			Command: "workbench.action.terminal.copySelectionAsHtml",
			When:    "terminalFocus",
		},
	}

	// Only add new bindings if not already present
	for _, newB := range newBindings {
		found := false
		for _, existingB := range existingBindings {
			if existingB.Key == newB.Key && existingB.Command == newB.Command && existingB.When == newB.When {
				found = true
				break
			}
		}
		if !found {
			existingBindings = append(existingBindings, newB)
		}
	}

	// Write updated list
	output, err := json.MarshalIndent(existingBindings, "", "    ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal keybindings: %v\n", err)
		return
	}

	if err := os.WriteFile(keybindingsPath, output, 0644); err != nil {
		fmt.Printf("❌ Failed to write keybindings.json: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully updated: %s\n", keybindingsPath)
}
