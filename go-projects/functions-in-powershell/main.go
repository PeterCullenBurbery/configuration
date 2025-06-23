package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func extractFunctionNames(node *yaml.Node, names *[]string) {
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			valNode := node.Content[i+1]
			if valNode.Kind == yaml.ScalarNode {
				line := strings.TrimSpace(valNode.Value)
				if strings.HasPrefix(line, "function") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						name := strings.TrimSpace(parts[1])
						name = strings.TrimRight(name, "{") // remove trailing brace if inline
						*names = append(*names, name)
					}
				}
			} else {
				extractFunctionNames(valNode, names)
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			extractFunctionNames(item, names)
		}
	}
}

func main() {
	variableName := "names.yaml"

	homeDir, _ := os.UserHomeDir()
	yamlPath := filepath.Join(homeDir, "Desktop", "GitHub-repositories", "configuration", "scripts.yaml")
	outputPath := filepath.Join(homeDir, "Desktop", "GitHub-repositories", "configuration", variableName)

	// Read scripts.yaml
	yamlBytes, err := os.ReadFile(yamlPath)
	if err != nil {
		panic(fmt.Errorf("❌ Failed to read YAML file: %w", err))
	}

	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		panic(fmt.Errorf("❌ Failed to parse YAML: %w", err))
	}

	var functionNames []string
	if len(root.Content) > 0 {
		extractFunctionNames(root.Content[0], &functionNames)
	}

	if len(functionNames) == 0 {
		fmt.Println("⚠️ No function names found.")
		return
	}

	// Build YAML list node
	namesNode := &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: []*yaml.Node{},
	}
	for _, name := range functionNames {
		namesNode.Content = append(namesNode.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: name,
			Tag:   "!!str",
		})
	}

	out := yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			namesNode,
		},
	}

	// If output file exists, back it up as .bak (overwrite if .bak exists)
	if _, err := os.Stat(outputPath); err == nil {
		backupPath := outputPath + ".bak"
		if err := os.Rename(outputPath, backupPath); err != nil {
			panic(fmt.Errorf("❌ Failed to backup existing file: %w", err))
		}
		fmt.Printf("📦 Existing file backed up to: %s\n", backupPath)
	}

	// Write to output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Errorf("❌ Failed to create output file: %w", err))
	}
	defer outFile.Close()

	encoder := yaml.NewEncoder(outFile)
	encoder.SetIndent(2)
	if err := encoder.Encode(&out); err != nil {
		panic(fmt.Errorf("❌ Failed to write YAML: %w", err))
	}
	encoder.Close()

	fmt.Printf("✅ Extracted %d function names to: %s\n", len(functionNames), outputPath)
}