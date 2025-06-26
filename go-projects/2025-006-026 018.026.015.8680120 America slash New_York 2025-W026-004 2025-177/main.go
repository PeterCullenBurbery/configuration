package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func mustRun(name string, args []string, dir string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to run %s %v: %v\n", name, args, err)
		os.Exit(1)
	}
}

func main() {
	basePath := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\go-projects\2025-006-026 018.026.044.1536530 America slash New_York 2025-W026-004 2025-177`

	// Ensure base path exists
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not create base path: %v\n", err)
		os.Exit(1)
	}

	// Step 1: run `go mod init multiply-test`
	fmt.Println("🚀 Initializing Go module...")
	mustRun("go", []string{"mod", "init", "multiply-test"}, basePath)

	// Step 2: create subfolder and multiply.go
	pkgDir := filepath.Join(basePath, "multiply-by-2718")
	if err := os.MkdirAll(pkgDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not create package directory: %v\n", err)
		os.Exit(1)
	}

	multiplyGo := `package multiplyby2718

// MultiplyBy2718 returns the input multiplied by 2718.
func MultiplyBy2718(n int) int {
	return n * 2718
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "multiply.go"), []byte(multiplyGo), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write multiply.go: %v\n", err)
		os.Exit(1)
	}

	// Step 3: create main.go
	mainGo := `package main

import (
	"fmt"
	"multiply-test/multiply-by-2718"
)

func main() {
	n := 2
	result := multiplyby2718.MultiplyBy2718(n)
	fmt.Printf("%d × 2718 = %d\n", n, result)
}
`
	if err := os.WriteFile(filepath.Join(basePath, "main.go"), []byte(mainGo), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write main.go: %v\n", err)
		os.Exit(1)
	}

	// Step 4: run go build
	fmt.Println("🔨 Building project...")
	mustRun("go", []string{"build", "-o", "multiply2718.exe"}, basePath)

	// Step 5: run the binary
	fmt.Println("🏃 Running the built executable...")
	mustRun(filepath.Join(basePath, "multiply2718.exe"), []string{}, basePath)
}