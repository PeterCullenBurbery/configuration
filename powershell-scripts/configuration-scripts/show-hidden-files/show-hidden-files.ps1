function Set-ShowHiddenFiles {
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

        # Set Hidden to 1 to show hidden files (DWORD)
        Set-ItemProperty -Path $regPath -Name Hidden -Value 1 -Type DWord
        Write-Host "✅ Hidden files will be shown (Hidden = 1)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply hidden files visibility."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to set show hidden files: $_"
    }
}