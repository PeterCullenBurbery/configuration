function Set-HideFileExtensions {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"

        # Ensure the key exists
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set HideFileExt to 1 to hide file extensions (DWORD)
        Set-ItemProperty -Path $regPath -Name HideFileExt -Value 1 -Type DWord
        Write-Host "✅ File extensions will be hidden (HideFileExt = 1)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply hiding of file extensions."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to hide file extensions: $_"
    }
}