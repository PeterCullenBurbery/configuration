function Set-LightMode {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize"

        # Ensure the key exists
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set light mode (1 = Light, DWORD)
        Set-ItemProperty -Path $regPath -Name AppsUseLightTheme   -Value 1 -Type DWord
        Set-ItemProperty -Path $regPath -Name SystemUsesLightTheme -Value 1 -Type DWord

        Write-Host "✅ Light mode registry values set (AppsUseLightTheme & SystemUsesLightTheme = 1)."

        if (-not $NoRestart) {
            # Restart Explorer so the new theme is read
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply Light Mode."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to set Light Mode: $_"
    }
}