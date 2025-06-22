function Set-StartMenuToCenter {
    $registryPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"
    $name = "TaskbarAl"

    try {
        # Set the registry value to 1 (center alignment)
        Set-ItemProperty -Path $registryPath -Name $name -Value 1 -Force

        # Restart Explorer (it will auto-restart)
        Stop-Process -Name explorer -Force

        Write-Host "✅ Start menu alignment set to center."
    } catch {
        Write-Warning "⚠️ Failed to set Start menu alignment: $_"
    }
}