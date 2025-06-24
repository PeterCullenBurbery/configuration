package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Get USERPROFILE environment variable
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		fmt.Println("❌ USERPROFILE environment variable not found.")
		return
	}

	// Construct the PowerShell 7 profile path
	profilePath := filepath.Join(userProfile, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	profileDir := filepath.Dir(profilePath)

	// Ensure the directory exists
	err := os.MkdirAll(profileDir, 0755)
	if err != nil {
		fmt.Printf("❌ Failed to create profile directory: %v\n", err)
		return
	}

	// PowerShell 7 profile content
	content := `# This is a comment

Import-Module MyModule

# --- Begin MyModule Logging Block PS7 ---
# --- Build Timestamp Filename ---

$now = Get-Date
$timestamp = $now.ToString("yyyy-0MM-0dd 0HH.0mm.0ss.fffffff")
$iana_tz = Get-IanaTimeZone
$iso_week_date = Get-IsoWeekDate -date $now
$iso_ordinal_date = Get-IsoOrdinalDate -date $now

$log_name = "$timestamp $iana_tz $iso_week_date $iso_ordinal_date"
$safe_log_name = $log_name -replace '/', ' slash '

$log_directory = "C:\terminal-logs\powershell-007-logs"
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
# --- End MyModule Logging Block PS007 ---

# Import the Chocolatey Profile that enables tab-completion for 'choco'
$ChocolateyProfile = "$env:ChocolateyInstall\helpers\chocolateyProfile.psm1"
if (Test-Path $ChocolateyProfile) {
    Import-Module "$ChocolateyProfile"
}`

	// Write to profile file
	err = os.WriteFile(profilePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("❌ Failed to write to PowerShell 7 profile: %v\n", err)
		return
	}

	fmt.Printf("✅ PowerShell 7 profile written to: %s\n", profilePath)
}
