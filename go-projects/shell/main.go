package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		// Read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Trim both \r and \n for Windows compatibility
		input = strings.TrimRight(input, "\r\n")

		// Skip empty input
		if input == "" {
			continue
		}

		// Execute the input
		if err := execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// ErrNoPath is returned when 'cd' is called without an argument.
var ErrNoPath = errors.New("path required")

func execInput(input string) error {
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	// Run using PowerShell 7 without loading the user profile
	cmd := exec.Command("pwsh", "-NoProfile", "-Command", input)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

