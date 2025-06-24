package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// --- Check if Windows PowerShell is available ---
	_, err := exec.LookPath("powershell.exe")
	if err != nil {
		fmt.Println("ℹ️ Windows PowerShell (powershell.exe) is not available. Skipping profile update.")
		return
	}

	// --- Get PowerShell 5 profile path ---
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", "Write-Output $PROFILE")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to get PowerShell 5 profile path: %v\n", err)
		return
	}

	profilePath := filepath.Clean(strings.TrimSpace(string(output)))

	// --- Ensure directory exists ---
	profileDir := filepath.Dir(profilePath)
	err = os.MkdirAll(profileDir, 0755)
	if err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		return
	}

	// --- Define PowerShell 5 profile content ---
	content := `# This is a comment

Import-Module MyModule

# --- Begin MyModule Logging Block PS5 ---
# --- Build Timestamp Filename ---

$now = Get-Date
$timestamp = $now.ToString("yyyy-0MM-0dd 0HH.0mm.0ss.fffffff")
$iana_tz = Get-IanaTimeZone
$iso_week_date = Get-IsoWeekDate -date $now
$iso_ordinal_date = Get-IsoOrdinalDate -date $now

$log_name = "$timestamp $iana_tz $iso_week_date $iso_ordinal_date"
$safe_log_name = $log_name -replace '/', ' slash '

$log_directory = "C:\terminal-logs\powershell-005-logs"
if (!(Test-Path $log_directory)) {
    New-Item -ItemType Directory -Path $log_directory | Out-Null
}

$log_file = Join-Path $log_directory "$safe_log_name.txt"

# --- Start Transcript ---
try {
    Start-Transcript -Path $log_file -Append -ErrorAction Stop
} catch {
    Write-Host "Transcript already running or failed to start."
}
# --- End MyModule Logging Block PS005 ---

# Import the Chocolatey Profile that enables tab-completion for 'choco'
$ChocolateyProfile = "$env:ChocolateyInstall\helpers\chocolateyProfile.psm1"
if (Test-Path $ChocolateyProfile) {
    Import-Module "$ChocolateyProfile"
}`

	// --- Write profile file ---
	err = os.WriteFile(profilePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("❌ Failed to write to PowerShell 5 profile: %v\n", err)
		return
	}

	fmt.Printf("✅ PowerShell 5 profile updated at: %s\n", profilePath)
}
