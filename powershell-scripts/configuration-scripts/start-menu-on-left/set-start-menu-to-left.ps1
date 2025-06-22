function Set-StartMenuToLeft {
    $registryPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"
    $name = "TaskbarAl"

    try {
        # Set the registry value to 0 (left alignment)
        Set-ItemProperty -Path $registryPath -Name $name -Value 0 -Force

        # Restart Explorer (Windows will auto-relaunch it)
        Stop-Process -Name explorer -Force

        Write-Host "✅ Start menu alignment set to left."
    } catch {
        Write-Warning "⚠️ Failed to set Start menu alignment: $_"
    }
}