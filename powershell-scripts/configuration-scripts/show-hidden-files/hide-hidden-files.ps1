function Set-HideHiddenFiles {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set Hidden to 2 to hide hidden files (DWORD)
        Set-ItemProperty -Path $regPath -Name Hidden -Value 2 -Type DWord
        Write-Host "✅ Hidden files will be hidden (Hidden = 2)."

        if (-not $NoRestart) {
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply hiding of hidden files."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to hide hidden files: $_"
    }
}