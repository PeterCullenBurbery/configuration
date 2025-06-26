# --- Start total timer ---
$global_start = Get-Date

# --- Define base path to local shared folder copy ---
$base_path = "\\vmware-host\Shared Folders\C\folder-for-shared-files-for-virtual-machines\GitHub-repositories\configuration"

# --- Unblock all files recursively to avoid execution prompts ---
Write-Host "🔓 Unblocking all files under: $base_path"
Get-ChildItem -Path $base_path -Recurse -File -ErrorAction SilentlyContinue | Unblock-File -ErrorAction SilentlyContinue

# --- Define important paths ---
$config_path      = Join-Path $base_path "config.yaml"
$install_path     = Join-Path $base_path "install.yaml"
$psm1_path        = Join-Path $base_path "output\MyModule.psm1"
$psd1_path        = Join-Path $base_path "output\MyModule.psd1"
$module_name      = [System.IO.Path]::GetFileNameWithoutExtension($psm1_path)
$go_command_cli   = Join-Path $base_path "go-projects\go-command-line\go-command-line.exe"

# --- Step 1: Copy PowerShell module ---
$copy_module_exe = Join-Path $base_path "go-projects\copy-module-to-c-powershell-modules\copy-module-to-c-powershell-modules.exe"

if ((Test-Path $copy_module_exe) -and (Test-Path $psm1_path) -and (Test-Path $psd1_path) -and (Test-Path $config_path)) {
    $start = Get-Date
    Write-Host "📦 Running copy-module-to-c-powershell-modules.exe..."
    Start-Process -FilePath $copy_module_exe `
        -ArgumentList "`"$psm1_path`"", "`"$psd1_path`"", "`"$config_path`"", "`"$module_name`"" `
        -Wait
    $end = Get-Date
    Write-Host "✅ Module copied in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
} else {
    Write-Error "❌ Missing files for module copy."
}

# --- Step 2: Add PowerShell module path via Go CLI ---
$add_ps_path_exe = Join-Path $base_path "go-projects\add-powershell-modules-to-path-by-executing-go-cli\add-powershell-modules-to-path-by-executing-go-cli.exe"

if ((Test-Path $add_ps_path_exe) -and (Test-Path $go_command_cli) -and (Test-Path $config_path)) {
    $start = Get-Date
    Write-Host "🧩 Running add-powershell-modules-to-path-by-executing-go-cli.exe..."
    Start-Process -FilePath $add_ps_path_exe `
        -ArgumentList "--cli", "`"$go_command_cli`"", "--yaml", "`"$config_path`"" `
        -Wait
    $end = Get-Date
    Write-Host "✅ PSModulePath updated in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
} else {
    Write-Error "❌ Missing CLI, Go tool, or config.yaml for PSModulePath step"
}

# --- Step 3: Run customize-file-explorer.exe ---
$explorer_exe = Join-Path $base_path "go-projects\customize-file-explorer\customize-file-explorer.exe"
$explorer_log = Join-Path $env:TEMP "$(Get-Date -Format 'yyyy-MM-dd_HH-mm-ss')_customize-file-explorer.log"

if ((Test-Path $explorer_exe) -and (Test-Path $config_path) -and (Test-Path $psm1_path)) {
    $start = Get-Date
    Write-Host "🖼️ Running customize-file-explorer.exe..."
    Start-Process -FilePath $explorer_exe `
        -ArgumentList "--config", "`"$config_path`"", "--module", "`"$psm1_path`"", "--log", "`"$explorer_log`"" `
        -Wait
    $end = Get-Date
    Write-Host "✅ customize-file-explorer completed in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
} else {
    Write-Error "❌ Missing files for customize-file-explorer"
}

# --- Step 4: Run install-things.exe ---
$install_exe = Join-Path $base_path "go-projects\install-things\install-things.exe"
$install_log = Join-Path $env:TEMP "$(Get-Date -Format 'yyyy-MM-dd_HH-mm-ss')_install-things.log"

if ((Test-Path $install_exe) -and (Test-Path $install_path) -and (Test-Path $psm1_path)) {
    $start = Get-Date
    Write-Host "🔧 Running install-things.exe..."
    Start-Process -FilePath $install_exe `
        -ArgumentList "--install", "`"$install_path`"", "--module", "`"$psm1_path`"", "--log", "`"$install_log`"" `
        -Wait
    $end = Get-Date
    Write-Host "✅ install-things completed in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
} else {
    Write-Error "❌ Missing files for install-things"
}

# --- Step 5: Run config tools (no arguments) ---
$config_tools = @(
    "configure-settings-for-windows-terminal\configure-settings-for-windows-terminal.exe",
    "configure-settings-for-vs-code\configure-settings-for-vs-code.exe",
    "configure-keyboard-shortcuts-for-vs-code\configure-keyboard-shortcuts-for-vs-code.exe"
)

foreach ($tool_rel in $config_tools) {
    $tool_path = Join-Path $base_path "go-projects\$tool_rel"
    if (Test-Path $tool_path) {
        $start = Get-Date
        Write-Host "⚙️ Running config tool: $tool_path"
        Start-Process -FilePath $tool_path -Wait
        $end = Get-Date
        Write-Host "✅ Finished $tool_path in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
    } else {
        Write-Warning "⚠️ Missing config tool: $tool_path"
    }
}

# --- Step 6: Update PowerShell profiles ---
$ps7_profile_exe = Join-Path $base_path "go-projects\powershell-007-profile\powershell-007-profile.exe"
$ps5_profile_exe = Join-Path $base_path "go-projects\powershell-005-profile\powershell-005-profile.exe"

if (Test-Path $ps5_profile_exe) {
    Write-Host "📜 Running PowerShell 5 profile updater..."
    Start-Process -FilePath $ps5_profile_exe -Wait
} else {
    Write-Warning "⚠️ PowerShell 5 profile updater not found"
}

if (Test-Path $ps7_profile_exe) {
    Write-Host "📜 Running PowerShell 7 profile updater..."
    Start-Process -FilePath $ps7_profile_exe -Wait
} else {
    Write-Warning "⚠️ PowerShell 7 profile updater not found"
}

# --- Step 7: Run enable-ssh.exe if SSH is enabled in config.yaml ---

$enable_ssh_exe = Join-Path $base_path "go-projects\enable-ssh\enable-ssh.exe"
$enable_ssh_log = Join-Path $env:TEMP "$(Get-Date -Format 'yyyy-MM-dd_HH-mm-ss')_enable-ssh.log"

if ((Test-Path $enable_ssh_exe) -and (Test-Path $config_path) -and (Test-Path $psm1_path)) {
    $start = Get-Date
    Write-Host "🔐 Running enable-ssh.exe..."

    Start-Process -FilePath $enable_ssh_exe `
        -ArgumentList "--yaml", "`"$config_path`"", "--module", "`"$psm1_path`"", "--log", "`"$enable_ssh_log`"" `
        -Wait

    $end = Get-Date
    Write-Host "✅ enable-ssh completed in $([math]::Round(($end - $start).TotalSeconds, 2)) seconds."
} else {
    Write-Warning "⚠️ Missing enable-ssh.exe, config.yaml, or PowerShell module."
}

# --- Print total execution time ---
$global_end = Get-Date
$total_seconds = [math]::Round(($global_end - $global_start).TotalSeconds, 2)
Write-Host "`n⏱️ Total execution time: $total_seconds seconds."

#install Visual Studio
# choco install visualstudio2022community --yes

# restart

# install Visual Studio C++ build tools
# Start-Process -FilePath "${env:ProgramFiles(x86)}\Microsoft Visual Studio\Installer\vs_installer.exe" -ArgumentList @(
#     'modify',
#     '--installPath', '"C:\Program Files\Microsoft Visual Studio\2022\Community"',
#     '--add', 'Microsoft.VisualStudio.Workload.NativeDesktop',
#     '--includeRecommended',
#     '--includeOptional',
#     '--passive',
#     '--norestart',
#     '--force'
# ) -Wait

# restart
