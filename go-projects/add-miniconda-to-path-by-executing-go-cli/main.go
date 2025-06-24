package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runGoCLI(goCLI, toolPath string) error {
	cmd := exec.Command(goCLI, "add-toPath", toolPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Define flags
	cliFlag := flag.String("cli", "", "Path to Go CLI executable")
	flag.Parse()

	// Fallback to positional argument
	var cliPath string
	if *cliFlag != "" {
		cliPath = *cliFlag
	} else if flag.NArg() == 1 {
		cliPath = flag.Arg(0)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  add_miniconda.exe --cli <go-cli>")
		fmt.Println("  add_miniconda.exe <go-cli>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Define tool paths
	baseDir := `C:\ProgramData\Miniconda3`
	pathsToAdd := []string{
		filepath.Join(baseDir, "python.exe"),
		filepath.Join(baseDir, "Scripts", "pip3.exe"),
	}

	// Run add-toPath for each
	for _, toolPath := range pathsToAdd {
		fmt.Printf("🚀 Running: %s add-toPath %s\n", cliPath, toolPath)

		if err := runGoCLI(cliPath, toolPath); err != nil {
			log.Fatalf("❌ Failed to add %s to PATH: %v\n", toolPath, err)
		}
	}

	fmt.Println("✅ Miniconda tools successfully added to PATH.")
}
