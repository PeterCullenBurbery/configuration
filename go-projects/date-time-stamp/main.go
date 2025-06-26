package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "java_hello_world")
	if err != nil {
		log.Fatalf("❌ Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Java file path
	javaFilePath := filepath.Join(tempDir, "HelloWorld.java")

	// Java source code
	javaCode := `public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World from Java!");
    }
}`

	// Write the Java file
	err = os.WriteFile(javaFilePath, []byte(javaCode), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write Java file: %v", err)
	}
	fmt.Println("📄 Java source written to:", javaFilePath)

	// Compile the Java source
	cmdCompile := exec.Command("javac", "HelloWorld.java")
	cmdCompile.Dir = tempDir
	cmdCompile.Stdout = os.Stdout
	cmdCompile.Stderr = os.Stderr
	if err := cmdCompile.Run(); err != nil {
		log.Fatalf("❌ Failed to compile Java code: %v", err)
	}
	fmt.Println("✅ Java compiled successfully.")

	// Run the compiled Java class
	cmdRun := exec.Command("java", "HelloWorld")
	cmdRun.Dir = tempDir
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	if err := cmdRun.Run(); err != nil {
		log.Fatalf("❌ Failed to run Java class: %v", err)
	}
	fmt.Println("🏁 Finished running Java program.")
}