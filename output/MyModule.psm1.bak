function Set-DarkMode {
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

        # Set dark mode (0 = Dark, DWORD)
        Set-ItemProperty -Path $regPath -Name AppsUseLightTheme   -Value 0 -Type DWord
        Set-ItemProperty -Path $regPath -Name SystemUsesLightTheme -Value 0 -Type DWord

        Write-Host "✅ Dark mode registry values set (AppsUseLightTheme & SystemUsesLightTheme = 0)."

        if (-not $NoRestart) {
            # Restart Explorer so the new theme is read
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply Dark Mode."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to set Dark Mode: $_"
    }
}

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

function Set-ShowFileExtensions {
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

        # Set HideFileExt to 0 to show file extensions (DWORD)
        Set-ItemProperty -Path $regPath -Name HideFileExt -Value 0 -Type DWord
        Write-Host "✅ File extensions will be visible (HideFileExt = 0)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply file extension visibility."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to set file extensions visibility: $_"
    }
}

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

function Set-HideSearchBox {
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

        # Set SearchboxTaskbarMode to 0 to hide the search box (DWORD)
        Set-ItemProperty -Path $regPath -Name SearchboxTaskbarMode -Value 0 -Type DWord
        Write-Host "✅ Search box will be hidden (SearchboxTaskbarMode = 0)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply hiding of search box."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to hide search box: $_"
    }
}

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

function Set-ShowSecondsInTaskbar {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"

        # Ensure the key exists (the path should exist by default)
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set ShowSecondsInSystemClock to 1 to show seconds (DWORD)
        Set-ItemProperty -Path $regPath -Name ShowSecondsInSystemClock -Value 1 -Type DWord
        Write-Host "✅ Taskbar clock will display seconds (ShowSecondsInSystemClock = 1)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply the seconds display."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to enable seconds in taskbar clock: $_"
    }
}

function Set-HideSecondsInTaskbar {
    [CmdletBinding()]
    param (
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )

    try {
        $regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"

        # Ensure the key exists (it normally does)
        if (-not (Test-Path $regPath)) {
            New-Item -Path $regPath -Force | Out-Null
        }

        # Set ShowSecondsInSystemClock to 0 to hide seconds (DWORD)
        Set-ItemProperty -Path $regPath -Name ShowSecondsInSystemClock -Value 0 -Type DWord
        Write-Host "✅ Taskbar clock will hide seconds (ShowSecondsInSystemClock = 0)."

        if (-not $NoRestart) {
            # Restart Explorer so the change takes effect immediately
            Stop-Process -Name explorer -Force
            Write-Host "🔁 Explorer restarted to apply the change."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to hide seconds in taskbar clock: $_"
    }
}

function Set-CustomShortDatePattern {
    [CmdletBinding()]
    param()

    $short_date_pattern = "yyyy-MM-dd-dddd"
    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Update the registry with the new short date pattern
        Set-ItemProperty -Path $reg_path -Name sShortDate -Value $short_date_pattern -ErrorAction Stop
        Write-Output "Short date pattern set to '$short_date_pattern'."

        # Define SendMessageTimeout P/Invoke only if not already defined
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "Failed to set short date pattern: $_"
    }
}

function Reset-ShortDatePattern {
    [CmdletBinding()]
    param(
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )
    try {
        # Determine the user's current culture name
        $cultureName = [System.Threading.Thread]::CurrentThread.CurrentCulture.Name

        # Create a CultureInfo without user overrides to get the default pattern
        $defaultCulture = New-Object System.Globalization.CultureInfo($cultureName, $false)
        $defaultPattern = $defaultCulture.DateTimeFormat.ShortDatePattern

        # Update registry
        $regPath = "HKCU:\Control Panel\International"
        Set-ItemProperty -Path $regPath -Name sShortDate -Value $defaultPattern -ErrorAction Stop

        # Broadcast setting change
        Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
using System;
using System.Runtime.InteropServices;
public static class NativeMethods {
    [DllImport("user32.dll", SetLastError = true)]
    public static extern IntPtr SendMessageTimeout(IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam, uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
}
'@ | Out-Null
        $HWND_BROADCAST = [IntPtr]0xFFFF
        $WM_SETTINGCHANGE = 0x001A
        [UIntPtr]$outPtr = [UIntPtr]::Zero
        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            [UIntPtr]0,
            "intl",
            0x0002,
            100,
            [ref]$outPtr
        ) | Out-Null

        Write-Host "✅ ShortDatePattern reset to default: '$defaultPattern'."

        if (-not $NoRestart) {
            Write-Host "🔁 Restarting Explorer to apply default date format..."
            Stop-Process -Name explorer -Force
            Write-Host "✅ Explorer restarted."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to reset short date pattern: $_"
    }
}

function Set-CustomLongDatePattern {
    [CmdletBinding()]
    param()

    $long_date_pattern = "yyyy-MM-dd-dddd"
    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Update the registry with the new long date pattern
        Set-ItemProperty -Path $reg_path -Name sLongDate -Value $long_date_pattern -ErrorAction Stop
        Write-Output "Long date pattern set to '$long_date_pattern'."

        # Define SendMessageTimeout only if not already defined
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "Failed to set long date pattern: $_"
    }
}

function Reset-LongDatePattern {
    [CmdletBinding()]
    param(
        [switch]$NoRestart  # If specified, do not restart Explorer automatically
    )
    try {
        # Determine current culture and default long date pattern without user overrides
        $cultureName = [System.Threading.Thread]::CurrentThread.CurrentCulture.Name
        $defaultCulture = New-Object System.Globalization.CultureInfo($cultureName, $false)
        $defaultPattern = $defaultCulture.DateTimeFormat.LongDatePattern

        # Update registry
        $regPath = "HKCU:\Control Panel\International"
        Set-ItemProperty -Path $regPath -Name sLongDate -Value $defaultPattern -ErrorAction Stop

        Write-Host "✅ LongDatePattern reset to default: '$defaultPattern'."

        if (-not $NoRestart) {
            Write-Host "🔁 Restarting Explorer to apply the change..."
            Stop-Process -Name explorer -Force
            Write-Host "✅ Explorer restarted."
        } else {
            Write-Host "ℹ️ Explorer restart skipped (use -NoRestart to prevent auto-restart)."
        }
    } catch {
        Write-Warning "❌ Failed to reset long date pattern: $_"
    }
}

function Set-CustomTimePattern {
    [CmdletBinding()]
    param()

    $time_format = "HH.mm.ss"   # Long time
    $short_time_format = "HH.mm.ss"  # Short time (with seconds)
    $time_separator = "."
    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Set long time, short time, and separator
        Set-ItemProperty -Path $reg_path -Name sTimeFormat -Value $time_format -ErrorAction Stop
        Set-ItemProperty -Path $reg_path -Name sShortTime -Value $short_time_format -ErrorAction Stop
        Set-ItemProperty -Path $reg_path -Name sTime -Value $time_separator -ErrorAction Stop

        Write-Output "✅ Time format set:"
        Write-Output "   Long time  (sTimeFormat): $time_format"
        Write-Output "   Short time (sShortTime) : $short_time_format"
        Write-Output "   Time separator (sTime)  : $time_separator"

        # Broadcast setting change
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply time settings."
    }
    catch {
        Write-Error "❌ Failed to set custom time pattern: $_"
    }
}

function Reset-TimePatternToDefault {
    [CmdletBinding()]
    param()

    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Get current culture without user overrides
        $culture = [System.Globalization.CultureInfo]::CurrentCulture
        $default_culture = New-Object System.Globalization.CultureInfo($culture.Name, $false)
        $default_long_time = $default_culture.DateTimeFormat.LongTimePattern
        $default_short_time = $default_culture.DateTimeFormat.ShortTimePattern
        $default_time_separator = $default_culture.DateTimeFormat.TimeSeparator

        # Reset time patterns and separator
        Set-ItemProperty -Path $reg_path -Name sTimeFormat -Value $default_long_time -ErrorAction Stop
        Set-ItemProperty -Path $reg_path -Name sShortTime -Value $default_short_time -ErrorAction Stop
        Set-ItemProperty -Path $reg_path -Name sTime -Value $default_time_separator -ErrorAction Stop

        Write-Output "✅ Time settings reset to system defaults:"
        Write-Output "   Long time  (sTimeFormat): $default_long_time"
        Write-Output "   Short time (sShortTime) : $default_short_time"
        Write-Output "   Time separator (sTime)  : $default_time_separator"

        # Broadcast setting change
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply default time settings."
    }
    catch {
        Write-Error "❌ Failed to reset time pattern to defaults: $_"
    }
}

function Set-24HourTimeFormat {
    [CmdletBinding()]
    param()

    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Set the 24-hour format flag
        Set-ItemProperty -Path $reg_path -Name iTime -Value "1" -ErrorAction Stop

        Write-Output "✅ Windows is now configured to use 24-hour time (iTime = 1)."

        # Optionally broadcast setting change to notify running apps
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xFFFF
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "❌ Failed to set 24-hour time format: $_"
    }
}

function Reset-12HourTimeFormat {
    [CmdletBinding()]
    param()

    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Set the 12-hour format flag
        Set-ItemProperty -Path $reg_path -Name iTime -Value "0" -ErrorAction Stop

        Write-Output "✅ Windows is now configured to use 12-hour time (iTime = 0)."

        # Optionally broadcast setting change to notify running apps
        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xFFFF
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "❌ Failed to reset to 12-hour time format: $_"
    }
}

function Set-FirstDayOfWeekMonday {
    [CmdletBinding()]
    param()

    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Set to Monday (0)
        Set-ItemProperty -Path $reg_path -Name iFirstDayOfWeek -Value "0" -ErrorAction Stop
        Write-Output "✅ First day of the week set to Monday (iFirstDayOfWeek = 0)."

        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xFFFF
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "❌ Failed to set first day of week: $_"
    }
}

function Set-FirstDayOfWeekSunday {

    [CmdletBinding()]
    param()

    $reg_path = "HKCU:\Control Panel\International"

    try {
        # Set to Sunday (6)
        Set-ItemProperty -Path $reg_path -Name iFirstDayOfWeek -Value "6" -ErrorAction Stop
        Write-Output "✅ First day of the week set to Sunday (iFirstDayOfWeek = 6)."

        if (-not ("Win32.NativeMethods" -as [type])) {
            Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true)]
public static extern int SendMessageTimeout(
    int hWnd, uint Msg, int wParam, string lParam, uint fuFlags, uint uTimeout, out int lpdwResult);
'@
        }

        $HWND_BROADCAST = 0xFFFF
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = 0

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            0,
            "Intl",
            $SMTO_ABORTIFHUNG,
            100,
            [ref]$result
        ) | Out-Null

        Write-Output "🔄 System broadcast completed to apply the setting."
    }
    catch {
        Write-Error "❌ Failed to set first day of week: $_"
    }
}

function Get-IanaTimeZone {
    $win_tz = (Get-TimeZone).Id
    $iana_tz = $null

    # Method 1: .NET TimeZoneInfo API (PowerShell 7.2+ / .NET 6+)
    if ([System.TimeZoneInfo].GetMethod("TryConvertWindowsIdToIanaId", [type[]]@([string],[string].MakeByRefType()))) {
        if ([System.TimeZoneInfo]::TryConvertWindowsIdToIanaId($win_tz, [ref] $iana_tz)) {
            return $iana_tz
        }
    }

    # Method 2: WinRT Calendar API (Windows 10+)
    try {
        return [Windows.Globalization.Calendar,Windows.Globalization,ContentType=WindowsRuntime]::new().GetTimeZone()
    } catch {}

    # Method 3: Parse TimeZoneMapping.xml
    $map_path = Join-Path $Env:WinDir 'Globalization\Time Zone\TimeZoneMapping.xml'
    if (Test-Path $map_path) {
        $map_xml = [xml](Get-Content $map_path)
        $node = $map_xml.TimeZoneMapping.MapTZ | Where-Object { $_.WinID -eq $win_tz -and $_.Default -eq "true" }
        if ($node) {
            return $node.TZID
        }
    }

    # Fallback to Windows ID
    return $win_tz
}

function Get-IsoWeekDate {
    param (
        [datetime]$date = (Get-Date)
    )

    if ([System.Type]::GetType("System.Globalization.ISOWeek")) {
        $iso_week = [System.Globalization.ISOWeek]::GetWeekOfYear($date)
        $iso_year = [System.Globalization.ISOWeek]::GetYear($date)
    } else {
        $iso_day = (([int]$date.DayOfWeek + 6) % 7) + 1
        $weekThursday = $date.AddDays(4 - $iso_day)
        $iso_year = $weekThursday.Year
        $iso_week = [System.Globalization.CultureInfo]::InvariantCulture.Calendar.GetWeekOfYear(
            $weekThursday,
            [System.Globalization.CalendarWeekRule]::FirstFourDayWeek,
            [System.DayOfWeek]::Monday
        )
    }

    $iso_day = (([int]$date.DayOfWeek + 6) % 7) + 1
    return "{0:0000}-W{1:000}-{2:000}" -f $iso_year, $iso_week, $iso_day
}

function Get-IsoOrdinalDate {
    [CmdletBinding()]
    param (
        [Parameter(ValueFromPipeline = $true)]
        [DateTime] $Date = (Get-Date)
    )

    process {
        # Format as YYYY-DDD (year and 3-digit day-of-year)
        $ordinal = "{0:yyyy}-{1:D3}" -f $Date, $Date.DayOfYear
        Write-Output $ordinal
    }
}

function prompt {
    $now = Get-Date

    # Use high-precision timestamp (7 fractional digits)
    $timestamp = $now.ToString("yyyy-0MM-0dd 0HH.0mm.0ss.fffffff")

    $iana_tz = Get-IanaTimeZone
    $iso_week_date = Get-IsoWeekDate -date $now
    $iso_ordinal_date = Get-IsoOrdinalDate -date $now

    # Print formatted info to screen (timestamp + tz + ISO week + ordinal)
    Write-Host "$timestamp $iana_tz $iso_week_date $iso_ordinal_date" -ForegroundColor Green

    return "$PWD> "
}

function Add-ToPath {

    param (
        [Parameter(Mandatory = $true)]
        [string]$PathToAdd
    )

    try {
        # Resolve path (file or folder)
        $resolvedPath = Resolve-Path -Path $PathToAdd -ErrorAction Stop

        # Determine if it's a file or folder
        if (Test-Path $resolvedPath.Path -PathType Leaf) {
            $targetPath = Split-Path -Path $resolvedPath.Path -Parent
        } else {
            $targetPath = $resolvedPath.Path
        }

        # Normalize
        $normalizedPath = $targetPath.TrimEnd('\')

        # Get current system PATH
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
        $pathEntries = $currentPath.Split(';') | ForEach-Object { $_.TrimEnd('\') }

        if ($pathEntries -contains $normalizedPath) {
            Write-Host "Path '$normalizedPath' is already in the system PATH."
            return
        }

        # Prepend and set new system PATH
        $newPath = "$normalizedPath;$currentPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
        Write-Host "Path '$normalizedPath' added to the TOP of system PATH."

        # Broadcast the environment change
        $signature = @'
[DllImport("user32.dll", SetLastError = true)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@

        Add-Type -MemberDefinition $signature -Namespace Win32 -Name NativeMethods

        $HWND_BROADCAST = [IntPtr]0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = [UIntPtr]::Zero

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            [UIntPtr]::Zero,
            "Environment",
            $SMTO_ABORTIFHUNG,
            5000,
            [ref]$result
        ) | Out-Null

        Write-Host "Environment update broadcast sent."
    }
    catch {
        Write-Error "Failed to add path: $_"
    }
}

function Remove-FromPath {
    param (
        [Parameter(Mandatory = $true)]
        [string]$PathToRemove
    )

    try {
        # Resolve to absolute path
        $resolvedPath = Resolve-Path -Path $PathToRemove -ErrorAction Stop

        # If it's a file, get the directory
        if (Test-Path $resolvedPath.Path -PathType Leaf) {
            $targetPath = Split-Path -Path $resolvedPath.Path -Parent
        } else {
            $targetPath = $resolvedPath.Path
        }

        # Normalize the path (remove trailing slashes)
        $normalizedPath = $targetPath.TrimEnd('\')

        # Get current system PATH
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
        $pathEntries = $currentPath.Split(';') | ForEach-Object { $_.TrimEnd('\') }

        if ($pathEntries -notcontains $normalizedPath) {
            Write-Host "Path '$normalizedPath' not found in system PATH."
            return
        }

        # Remove the path
        $updatedPath = ($pathEntries | Where-Object { $_ -ne $normalizedPath }) -join ';'
        [Environment]::SetEnvironmentVariable("Path", $updatedPath, "Machine")
        Write-Host "Path '$normalizedPath' removed from system PATH."

        # Broadcast the environment variable change
        $signature = @'
[DllImport("user32.dll", SetLastError = true)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@

        Add-Type -MemberDefinition $signature -Namespace Win32 -Name NativeMethods

        $HWND_BROADCAST = [IntPtr]0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = [UIntPtr]::Zero

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            [UIntPtr]::Zero,
            "Environment",
            $SMTO_ABORTIFHUNG,
            5000,
            [ref]$result
        ) | Out-Null

        Write-Host "Environment update broadcast sent."
    }
    catch {
        Write-Error "Failed to remove path: $_"
    }
}

function Get-SystemPath {
    [CmdletBinding()]
    param ()

    $path = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
    return $path
}

function Add-ToPSModulePath {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$Directory
    )

    try {
        # Resolve and normalize the path
        $resolvedPath = (Resolve-Path -Path $Directory -ErrorAction Stop).Path.TrimEnd('\')

        # Get current system PSModulePath
        $currentSystemPath = [Environment]::GetEnvironmentVariable("PSModulePath", "Machine")
        $pathEntries = $currentSystemPath -split ';' | ForEach-Object { $_.TrimEnd('\') }

        if ($pathEntries -contains $resolvedPath) {
            Write-Host "Path already in system PSModulePath: $resolvedPath"
        } else {
            $newPath = "$currentSystemPath;$resolvedPath"
            [Environment]::SetEnvironmentVariable("PSModulePath", $newPath, "Machine")
            Write-Host "Added to system PSModulePath: $resolvedPath"

            # Optional: Broadcast the environment change to running processes
            $signature = @'
[DllImport("user32.dll", SetLastError = true)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@

            Add-Type -MemberDefinition $signature -Namespace Win32 -Name NativeMethods
            $HWND_BROADCAST = [IntPtr]0xffff
            $WM_SETTINGCHANGE = 0x001A
            $SMTO_ABORTIFHUNG = 0x0002
            $result = [UIntPtr]::Zero

            [Win32.NativeMethods]::SendMessageTimeout(
                $HWND_BROADCAST,
                $WM_SETTINGCHANGE,
                [UIntPtr]::Zero,
                "Environment",
                $SMTO_ABORTIFHUNG,
                5000,
                [ref]$result
            ) | Out-Null

            Write-Host "Environment update broadcast sent."
        }
    } catch {
        Write-Error "Failed to add to PSModulePath: $_"
    }
}

function Remove-FromPSModulePath {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$Directory
    )

    try {
        # Resolve and normalize the path
        $resolvedPath = (Resolve-Path -Path $Directory -ErrorAction Stop).Path.TrimEnd('\')

        # Get current system PSModulePath
        $currentSystemPath = [Environment]::GetEnvironmentVariable("PSModulePath", "Machine")
        $pathEntries = $currentSystemPath -split ';' | ForEach-Object { $_.TrimEnd('\') }

        if ($pathEntries -notcontains $resolvedPath) {
            Write-Host "Path not found in system PSModulePath: $resolvedPath"
            return
        }

        # Remove the path
        $updatedPath = ($pathEntries | Where-Object { $_ -ne $resolvedPath }) -join ';'
        [Environment]::SetEnvironmentVariable("PSModulePath", $updatedPath, "Machine")
        Write-Host "Removed from system PSModulePath: $resolvedPath"

        # Broadcast the environment change
        $signature = @'
[DllImport("user32.dll", SetLastError = true)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@

        Add-Type -MemberDefinition $signature -Namespace Win32 -Name NativeMethods
        $HWND_BROADCAST = [IntPtr]0xffff
        $WM_SETTINGCHANGE = 0x001A
        $SMTO_ABORTIFHUNG = 0x0002
        $result = [UIntPtr]::Zero

        [Win32.NativeMethods]::SendMessageTimeout(
            $HWND_BROADCAST,
            $WM_SETTINGCHANGE,
            [UIntPtr]::Zero,
            "Environment",
            $SMTO_ABORTIFHUNG,
            5000,
            [ref]$result
        ) | Out-Null

        Write-Host "Environment update broadcast sent."
    }
    catch {
        Write-Error "Failed to remove from PSModulePath: $_"
    }
}

function Add-DomainAdminUser {

    param (

        [Parameter(Mandatory = $true)]
        [string]$username

    )

    $password = ConvertTo-SecureString "f" -AsPlainText -Force

    # Create the domain user
    New-ADUser `
        -Name $username `
        -SamAccountName $username `
        -UserPrincipalName "$username@corp.strength.local" `
        -AccountPassword $password `
        -Enabled $true `
        -PasswordNeverExpires $true `
        -Path "CN=Users,DC=corp,DC=strength,DC=local"

    # Add to Domain Admins
    Add-ADGroupMember -Identity "Domain Admins" -Members $username

    # Add to Remote Desktop Users
    Add-ADGroupMember -Identity "Remote Desktop Users" -Members $username

    Write-Host "User '$username' created and added to Domain Admins and Remote Desktop Users."
}

function Test-Is64Bit {
    if ([Environment]::Is64BitOperatingSystem) {
        Write-Host "Your computer is 64-bit."
    } else {
        Write-Host "Your computer is 32-bit."
    }
}

function Install-PowerShell-7 {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of PowerShell 7..."

    $arguments = @(
        "install"
        "--id", "Microsoft.PowerShell"
        "--source", "winget"
        "--scope", "machine"
        "--silent"
        "--accept-package-agreements"
        "--accept-source-agreements"
    )

    try {
        Start-Process -FilePath "winget" -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ PowerShell 7 installed successfully."
    } catch {
        Write-Error "❌ Failed to install PowerShell 7. Error: $_"
    }
}

function Install-VSCode {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of Visual Studio Code..."

    $arguments = @(
        "install"
        "-e"
        "--id", "Microsoft.VisualStudioCode"
        "--scope", "machine"
        "--silent"
        "--accept-package-agreements"
        "--accept-source-agreements"
    )

    try {
        Start-Process -FilePath "winget" -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ Visual Studio Code installed successfully."
    } catch {
        Write-Error "❌ Failed to install Visual Studio Code. Error: $_"
    }
}

function Install-7Zip {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of 7-Zip..."

    $arguments = @(
        "install"
        "-e"
        "--id", "7zip.7zip"
        "--scope", "machine"
        "--silent"
        "--accept-package-agreements"
        "--accept-source-agreements"
    )

    try {
        Start-Process -FilePath "winget" -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ 7-Zip installed successfully."
    } catch {
        Write-Error "❌ Failed to install 7-Zip. Error: $_"
    }
}

function Install-Voidtools-Everything {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of Voidtools Everything..."

    $arguments = @(
        "install"
        "-e"
        "--id", "voidtools.Everything"
        "--scope", "machine"
        "--silent"
        "--accept-package-agreements"
        "--accept-source-agreements"
    )

    try {
        Start-Process -FilePath "winget" -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ Voidtools Everything installed successfully."
    } catch {
        Write-Error "❌ Failed to install Voidtools Everything. Error: $_"
    }
}

function Install-WinSCP {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of WinSCP..."

    $arguments = @(
        "install"
        "-e"
        "--id", "WinSCP.WinSCP"
        "--scope", "machine"
        "--silent"
        "--accept-package-agreements"
        "--accept-source-agreements"
    )

    try {
        Start-Process -FilePath "winget" -ArgumentList $arguments -Wait -NoNewWindow
        Write-Host "✅ WinSCP installed successfully."
    } catch {
        Write-Error "❌ Failed to install WinSCP. Error: $_"
    }
}

function Install-MobaXterm {

    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of MobaXterm..."

    # Try to resolve choco path
    $chocoPath = Get-Command choco -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source

    if (-not $chocoPath) {
        $defaultChoco = "C:\ProgramData\chocolatey\bin\choco.exe"
        if (Test-Path $defaultChoco) {
            $chocoPath = $defaultChoco
        } else {
            Write-Error "❌ Chocolatey not found. Please install Chocolatey first."
            return
        }
    }

    $arguments = @("install", "mobaxterm", "--yes")

    try {
        Start-Process -FilePath $chocoPath -ArgumentList $arguments -Wait -NoNewWindow

        # Confirm installation
        $isInstalled = & $chocoPath list --local-only | Select-String -Pattern '^mobaxterm'

        if ($isInstalled) {
            Write-Host "✅ MobaXterm installed successfully."
        } else {
            Write-Warning "⚠️ MobaXterm install completed, but it may not be installed correctly."
        }

    } catch {
        Write-Error "❌ Failed to install MobaXterm. Error: $_"
    }
}

function Install-Go {

    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of Go..."

    # Resolve choco path
    $chocoPath = Get-Command choco -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source

    if (-not $chocoPath) {
        $defaultChoco = "C:\ProgramData\chocolatey\bin\choco.exe"
        if (Test-Path $defaultChoco) {
            $chocoPath = $defaultChoco
        } else {
            Write-Error "❌ Chocolatey not found. Please install Chocolatey first."
            return
        }
    }

    $arguments = @("install", "golang", "--yes")

    try {
        Start-Process -FilePath $chocoPath -ArgumentList $arguments -Wait -NoNewWindow

        # Verify installation
        $isInstalled = & $chocoPath list --local-only | Select-String -Pattern '^golang'

        if ($isInstalled) {
            Write-Host "✅ Go installed successfully."
        } else {
            Write-Warning "⚠️ Installation completed, but Go may not be fully installed."
        }

    } catch {
        Write-Error "❌ Failed to install Go. Error: $_"
    }
}

function Install-NotepadPP {

    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of Notepad++..."

    # Resolve Chocolatey path
    $chocoPath = Get-Command choco -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source

    if (-not $chocoPath) {
        $defaultChoco = "C:\ProgramData\chocolatey\bin\choco.exe"
        if (Test-Path $defaultChoco) {
            $chocoPath = $defaultChoco
        } else {
            Write-Error "❌ Chocolatey not found. Please install Chocolatey first."
            return
        }
    }

    $arguments = @("install", "notepadplusplus", "--yes")

    try {
        Start-Process -FilePath $chocoPath -ArgumentList $arguments -Wait -NoNewWindow

        # Verify installation
        $isInstalled = & $chocoPath list --local-only | Select-String -Pattern '^notepadplusplus'

        if ($isInstalled) {
            Write-Host "✅ Notepad++ installed successfully."
        } else {
            Write-Warning "⚠️ Install command ran, but Notepad++ may not be fully installed."
        }

    } catch {
        Write-Error "❌ Failed to install Notepad++. Error: $_"
    }
}

function Install-SQLiteBrowser {

    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of DB Browser for SQLite..."

    # Resolve Chocolatey path
    $chocoPath = Get-Command choco -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source

    if (-not $chocoPath) {
        $defaultChoco = "C:\ProgramData\chocolatey\bin\choco.exe"
        if (Test-Path $defaultChoco) {
            $chocoPath = $defaultChoco
        } else {
            Write-Error "❌ Chocolatey not found. Please install Chocolatey first."
            return
        }
    }

    $arguments = @("install", "sqlitebrowser", "--yes")

    try {
        Start-Process -FilePath $chocoPath -ArgumentList $arguments -Wait -NoNewWindow

        # Verify installation
        $isInstalled = & $chocoPath list --local-only | Select-String -Pattern '^sqlitebrowser'

        if ($isInstalled) {
            Write-Host "✅ DB Browser for SQLite installed successfully."
        } else {
            Write-Warning "⚠️ Install completed, but DB Browser for SQLite may not be fully installed."
        }

    } catch {
        Write-Error "❌ Failed to install DB Browser for SQLite. Error: $_"
    }
}

function Install-Java {

    [CmdletBinding()]
    param (
        [string]$PackageName = "temurin21"
    )

    Write-Host "🚀 Starting installation of Java package: $PackageName..."

    # Resolve Chocolatey path
    $chocoPath = Get-Command choco -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source
    if (-not $chocoPath) {
        $defaultChoco = "C:\ProgramData\chocolatey\bin\choco.exe"
        if (Test-Path $defaultChoco) {
            $chocoPath = $defaultChoco
        } else {
            Write-Error "❌ Chocolatey not found. Please install Chocolatey first."
            return
        }
    }

    $arguments = @("install", $PackageName, "--yes")

    try {
        Start-Process -FilePath $chocoPath -ArgumentList $arguments -Wait -NoNewWindow

        # Confirm installation
        $isInstalled = & $chocoPath list --local-only | Select-String -Pattern "^\Q$PackageName\E"
        if ($isInstalled) {
            Write-Host "✅ Java ($PackageName) installed successfully."
        } else {
            Write-Warning "⚠️ Install completed, but $PackageName may not be fully installed."
        }

        # Set JAVA_HOME
        $javaHomePath = $null

        if ($PackageName -ieq "temurin21") {
            $javaHomePath = "C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot"
        } else {
            $jdkDir = Get-ChildItem "C:\Program Files\Eclipse Adoptium\" -Directory |
                      Where-Object { $_.Name -like "jdk*" } |
                      Sort-Object LastWriteTime -Descending |
                      Select-Object -First 1
            if ($jdkDir) {
                $javaHomePath = $jdkDir.FullName
            }
        }

        if ($javaHomePath -and (Test-Path $javaHomePath)) {
            [Environment]::SetEnvironmentVariable("JAVA_HOME", $javaHomePath, [System.EnvironmentVariableTarget]::Machine)
            Write-Host "🌱 JAVA_HOME auto-set to: $javaHomePath"
        } else {
            Write-Warning "⚠️ Could not determine JAVA_HOME path. You may need to set it manually."
        }

    } catch {
        Write-Error "❌ Failed to install Java ($PackageName). Error: $_"
    }
}

function Install-CherryTree {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $true)]
        [string]$log,

        [Parameter(Mandatory = $true)]
        [string]$installPath
    )

    # Logging helper
    function Write-Log {
        param ([string]$message)
        $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        "$timestamp`t$message" | Out-File -FilePath $log -Append -Encoding UTF8
    }

    Write-Host "🚀 Starting CherryTree installation..."
    Write-Host "📝 Log path: $log"
    Write-Host "📁 Install path: $installPath"

    # Ensure directories exist
    $logDir = Split-Path $log -Parent
    if (-not (Test-Path $logDir)) {
        New-Item -ItemType Directory -Path $logDir -Force | Out-Null
    }

    $installDirParent = Split-Path $installPath -Parent
    if (-not (Test-Path $installDirParent)) {
        New-Item -ItemType Directory -Path $installDirParent -Force | Out-Null
    }

    # Use dynamic path to installer
    $installerName = "cherrytree_1.5.0.0_win64_setup.exe"
    $installer = Join-Path $installPath $installerName
    if (-not (Test-Path $installer)) {
        Write-Log "❌ Installer not found at $installer"
        Write-Error "❌ Installer not found at $installer"
        return
    }

    # Install arguments
    $arguments = @(
        "/VERYSILENT"
        "/SUPPRESSMSGBOXES"
        "/NORESTART"
        "/SP-"
        "/DIR=$installPath"
        "/LOG=$log"
    )

    $start = Get-Date
    Write-Log "🚀 Install started"
    Write-Host "⏱️ Start: $start"

    try {
        Start-Process -FilePath $installer -ArgumentList $arguments -Wait -NoNewWindow

        $end = Get-Date
        $duration = $end - $start

        Write-Log "✅ Install completed"
        Write-Log "⏱️ Start: $start"
        Write-Log "✅ End:   $end"
        Write-Log "🧮 Duration: $($duration.ToString())"

        Write-Host "✅ End:   $end"
        Write-Host "🧮 Duration: $($duration.ToString())"
    } catch {
        $end = Get-Date
        $duration = $end - $start

        Write-Log "❌ Install failed: $_"
        Write-Log "⏱️ Start: $start"
        Write-Log "❌ End:   $end"
        Write-Log "🧮 Duration: $($duration.ToString())"

        Write-Error "❌ Installation failed"
        Write-Host "⏱️ Start: $start"
        Write-Host "❌ End:   $end"
        Write-Host "🧮 Duration: $($duration.ToString())"
    }
}

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

function Install-Choco {
    [CmdletBinding()]
    param ()

    Write-Host "🚀 Starting installation of Chocolatey..."

    $installScript = 'https://community.chocolatey.org/install.ps1'

    try {
        Set-ExecutionPolicy Bypass -Scope Process -Force

        # Secure protocol
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072

        # Run install script
        iex ((New-Object System.Net.WebClient).DownloadString($installScript))

        Write-Host "✅ Chocolatey installed successfully."
    } catch {
        Write-Error "❌ Failed to install Chocolatey. Error: $_"
    }
}