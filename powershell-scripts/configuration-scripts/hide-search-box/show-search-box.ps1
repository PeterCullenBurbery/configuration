function Set-ShowSearchBox {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Search"
        # Ensure the key exists
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set SearchboxTaskbarMode to 2 to show the full search box (DWORD)
        Set-ItemProperty -Path $regPath -Name SearchboxTaskbarMode -Value 2 -Type DWord
        Write-Host "✅ Search box will be shown (SearchboxTaskbarMode = 2)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply showing of search box."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to show search box: $_"
    }
}