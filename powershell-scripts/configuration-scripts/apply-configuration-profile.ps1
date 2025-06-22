param (
    [string]$config_path = "$PSScriptRoot\config.yaml"
)

# --- Ensure NuGet provider is installed silently ---
if (-not (Get-PackageProvider -Name NuGet -ErrorAction SilentlyContinue)) {
    Write-Host "📦 Installing NuGet provider..."
    try {
        Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force -Scope AllUsers
    } catch {
        Write-Error "❌ Failed to install NuGet provider. Please run this script as Administrator."
        exit 1
    }
}

# --- Ensure powershell-yaml is installed for all users ---
if (-not (Get-Module -ListAvailable -Name powershell-yaml)) {
    Write-Host "📦 Installing 'powershell-yaml' module for all users..."
    try {
        Install-Module -Name powershell-yaml -Force -Scope AllUsers -ErrorAction Stop
    } catch {
        Write-Error "❌ Failed to install 'powershell-yaml'. Please run this script as Administrator."
        exit 1
    }
}

# --- Import powershell-yaml ---
Import-Module powershell-yaml -ErrorAction Stop

# --- Load YAML config ---
if (-Not (Test-Path $config_path)) {
    Write-Error "❌ Config file not found at $config_path"
    exit 1
}
$config = Get-Content $config_path | ConvertFrom-Yaml

# --- Apply configuration ---
Write-Host "🛠️ Applying configuration from $config_path..."

$base = $PSScriptRoot

# Dark mode
if ($config.configuration_profile.dark_mode -eq $true) {
    & "$base\dark-mode\set-dark-mode.ps1"
} else {
    & "$base\dark-mode\set-light-mode.ps1"
}

# Search box
switch ($config.configuration_profile.search_box) {
    "hidden" { & "$base\hide-search-box\hide-search-box.ps1" }
    "shown"  { & "$base\hide-search-box\show-search-box.ps1" }
}

# File extensions
switch ($config.configuration_profile.file_extensions) {
    "shown"  { & "$base\show-file-extensions\show-file-extensions.ps1" }
    "hidden" { & "$base\show-file-extensions\hide-file-extensions.ps1" }
}

# Hidden files
switch ($config.configuration_profile.hidden_files) {
    "shown"  { & "$base\show-hidden-files\show-hidden-files.ps1" }
    "hidden" { & "$base\show-hidden-files\hide-hidden-files.ps1" }
}

# Start menu alignment
switch ($config.configuration_profile.start_menu_alignment) {
    "left"   { & "$base\start-menu-on-left\set-start-menu-to-left.ps1" }
    "center" { & "$base\start-menu-on-left\set-start-menu-to-center.ps1" }
}

Write-Host "✅ Configuration applied."
