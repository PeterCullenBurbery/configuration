package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var powershellFunctions = map[string]string{
	"add-topath": `
function Add-ToPath {
    param ([string]$PathToAdd)
    try {
        $resolvedPath = Resolve-Path -Path $PathToAdd -ErrorAction Stop
        $targetPath = if (Test-Path $resolvedPath.Path -PathType Leaf) {
            Split-Path -Path $resolvedPath.Path -Parent
        } else {
            $resolvedPath.Path
        }
        $normalizedPath = $targetPath.TrimEnd('\')
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
        $pathEntries = $currentPath.Split(';') | ForEach-Object { $_.TrimEnd('\') }
        if ($pathEntries -contains $normalizedPath) {
            Write-Host "Path '$normalizedPath' is already in the system PATH."
            return
        }
        $newPath = "$normalizedPath;$currentPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
        Write-Host "Path '$normalizedPath' added to the TOP of system PATH."
        Broadcast-EnvChange
    } catch {
        Write-Error "Failed to add path: $_"
    }
}
`,
	"remove-frompath": `
function Remove-FromPath {
    param ([string]$PathToRemove)
    try {
        $resolvedPath = Resolve-Path -Path $PathToRemove -ErrorAction Stop
        $targetPath = if (Test-Path $resolvedPath.Path -PathType Leaf) {
            Split-Path -Path $resolvedPath.Path -Parent
        } else {
            $resolvedPath.Path
        }
        $normalizedPath = $targetPath.TrimEnd('\')
        $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
        $pathEntries = $currentPath.Split(';') | ForEach-Object { $_.TrimEnd('\') }
        if ($pathEntries -notcontains $normalizedPath) {
            Write-Host "Path '$normalizedPath' not found in system PATH."
            return
        }
        $updatedPath = ($pathEntries | Where-Object { $_ -ne $normalizedPath }) -join ';'
        [Environment]::SetEnvironmentVariable("Path", $updatedPath, "Machine")
        Write-Host "Path '$normalizedPath' removed from system PATH."
        Broadcast-EnvChange
    } catch {
        Write-Error "Failed to remove path: $_"
    }
}
`,
	"add-topsmodulepath": `
function Add-ToPSModulePath {
    param ([string]$Directory)
    try {
        $resolvedPath = (Resolve-Path -Path $Directory -ErrorAction Stop).Path.TrimEnd('\')
        $currentSystemPath = [Environment]::GetEnvironmentVariable("PSModulePath", "Machine")
        $pathEntries = $currentSystemPath -split ';' | ForEach-Object { $_.TrimEnd('\') }
        if ($pathEntries -contains $resolvedPath) {
            Write-Host "Path already in system PSModulePath: $resolvedPath"
        } else {
            $newPath = "$currentSystemPath;$resolvedPath"
            [Environment]::SetEnvironmentVariable("PSModulePath", $newPath, "Machine")
            Write-Host "Added to system PSModulePath: $resolvedPath"
            Broadcast-EnvChange
        }
    } catch {
        Write-Error "Failed to add to PSModulePath: $_"
    }
}
`,
	"remove-frompsmodulepath": `
function Remove-FromPSModulePath {
    param ([string]$Directory)
    try {
        $resolvedPath = (Resolve-Path -Path $Directory -ErrorAction Stop).Path.TrimEnd('\')
        $currentSystemPath = [Environment]::GetEnvironmentVariable("PSModulePath", "Machine")
        $pathEntries = $currentSystemPath -split ';' | ForEach-Object { $_.TrimEnd('\') }
        if ($pathEntries -notcontains $resolvedPath) {
            Write-Host "Path not found in system PSModulePath: $resolvedPath"
            return
        }
        $updatedPath = ($pathEntries | Where-Object { $_ -ne $resolvedPath }) -join ';'
        [Environment]::SetEnvironmentVariable("PSModulePath", $updatedPath, "Machine")
        Write-Host "Removed from system PSModulePath: $resolvedPath"
        Broadcast-EnvChange
    } catch {
        Write-Error "Failed to remove from PSModulePath: $_"
    }
}
`,
	"broadcast-envchange": `
function Broadcast-EnvChange {
    $signature = @"
    [DllImport("user32.dll", SetLastError = true)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
        uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
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
}
`,
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:")
		fmt.Println("  1234.exe add-toPath <path>")
		fmt.Println("  1234.exe remove-fromPath <path>")
		fmt.Println("  1234.exe add-toPSModulePath <path>")
		fmt.Println("  1234.exe remove-fromPSModulePath <path>")
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])
	pathArg := os.Args[2]

	scriptFunc, found := powershellFunctions[command]
	if !found {
		fmt.Println("❌ Unsupported command:", command)
		os.Exit(1)
	}

	fullScript := powershellFunctions["broadcast-envchange"] + scriptFunc + fmt.Sprintf("\n%s \"%s\"", getFunctionCall(command), pathArg)

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", fullScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ PowerShell execution failed:", err)
		os.Exit(1)
	}
}

func getFunctionCall(command string) string {
	switch command {
	case "add-topath":
		return "Add-ToPath"
	case "remove-frompath":
		return "Remove-FromPath"
	case "add-topsmodulepath":
		return "Add-ToPSModulePath"
	case "remove-frompsmodulepath":
		return "Remove-FromPSModulePath"
	default:
		return ""
	}
}
