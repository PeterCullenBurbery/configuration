function Install-Miniconda {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$InstallerPath
    )

    # Define installation paths
    $installDir = "C:\ProgramData\Miniconda3"
    $condaExe = Join-Path $installDir "Scripts\conda.exe"
    $pythonExe = Join-Path $installDir "python.exe"

    # Check if installer exists
    if (-not (Test-Path -Path $InstallerPath)) {
        Write-Error "❌ Installer not found at: $InstallerPath"
        return
    }

    Write-Host "📦 Installing Miniconda from: $InstallerPath"

    # Define install arguments
    $arguments = @(
        "/S",                                # Silent install
        "/InstallationType=AllUsers",        # System-wide
        "/AddToPath=1",                      # Add to PATH
        "/RegisterPython=1",                 # Set as system Python
        "/D=$installDir"                     # Install location (must be last)
    )

    try {
        # Run installer
        Start-Process -FilePath $InstallerPath -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ Miniconda installed successfully."

        # --- Verification ---
        Write-Host "`n✅ Miniconda installed to: $installDir"

        if (Test-Path $pythonExe) {
            Write-Host "🐍 Python version:"
            & $pythonExe --version
        } else {
            Write-Warning "⚠️ Python not found at expected path: $pythonExe"
        }

        if (Test-Path $condaExe) {
            Write-Host "📦 Conda version:"
            & $condaExe --version

            # Clear Conda cache
            & $condaExe clean --all --yes
            Write-Host "🧹 Conda cache cleaned."
        } else {
            Write-Warning "⚠️ Conda not found at expected path: $condaExe"
        }
    } catch {
        Write-Error "❌ Installation failed: $_"
    }
}