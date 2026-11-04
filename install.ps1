# jvm Windows PowerShell Installer (irm https://.../install.ps1 | iex)

$ErrorActionPreference = "Stop"

$jvmDir = if ($env:JVM_DIR) { $env:JVM_DIR } else { "$HOME\.jvm" }

Write-Host "Installing jvm to $jvmDir..." -ForegroundColor Cyan

$binDir = "$jvmDir\bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null
New-Item -ItemType Directory -Force -Path "$jvmDir\versions" | Out-Null

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

if (Test-Path "$scriptDir\bin\jvm.exe") {
    Copy-Item "$scriptDir\bin\jvm.exe" -Destination "$binDir\jvm.exe" -Force
    Copy-Item "$scriptDir\jvm.ps1" -Destination "$jvmDir\jvm.ps1" -Force
    if (Test-Path "$scriptDir\jvm.cmd") { Copy-Item "$scriptDir\jvm.cmd" -Destination "$binDir\jvm.cmd" -Force }
} elseif (Test-Path ".\bin\jvm.exe") {
    Copy-Item ".\bin\jvm.exe" -Destination "$binDir\jvm.exe" -Force
    Copy-Item ".\jvm.ps1" -Destination "$jvmDir\jvm.ps1" -Force
    if (Test-Path ".\jvm.cmd") { Copy-Item ".\jvm.cmd" -Destination "$binDir\jvm.cmd" -Force }
} else {
    # Download latest release binary from GitHub
    $repo = "unstoppable-tanmay/jvm"
    $arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    Write-Host "Fetching latest release from GitHub ($arch)..." -ForegroundColor Cyan
    try {
        $releaseUrl = "https://api.github.com/repos/$repo/releases/latest"
        $release = Invoke-RestMethod -Uri $releaseUrl -Headers @{ "User-Agent" = "jvm-installer" }
        $tag = $release.tag_name
        $version = $tag.TrimStart('v')
        $zipUrl = "https://github.com/$repo/releases/download/$tag/jvm_${version}_windows_${arch}.zip"
        $tempZip = Join-Path $env:TEMP "jvm_${version}.zip"

        Write-Host "Downloading $zipUrl..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $zipUrl -OutFile $tempZip
        Expand-Archive -Path $tempZip -DestinationPath $binDir -Force
        Remove-Item $tempZip -Force
    } catch {
        Write-Warning "Could not fetch pre-built binary: $_"
    }

    # Download shell scripts
    try {
        Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$repo/main/jvm.ps1" -OutFile "$jvmDir\jvm.ps1"
        Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$repo/main/jvm.cmd" -OutFile "$binDir\jvm.cmd"
    } catch {
        Write-Warning "Could not download wrapper scripts: $_"
    }
}

# Add $binDir to User Environment PATH permanently so both PowerShell and CMD recognize jvm and java shims
$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if ($userPath -notlike "*$binDir*") {
    $newUserPath = if ($userPath) { "$binDir;$userPath" } else { $binDir }
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, [EnvironmentVariableTarget]::User)
    Write-Host "✓ Added $binDir to User PATH environment variable (works in CMD, PowerShell & non-interactive)" -ForegroundColor Green
}

$profilePath = $PROFILE
if (-not (Test-Path (Split-Path $profilePath))) {
    New-Item -ItemType Directory -Force -Path (Split-Path $profilePath) | Out-Null
}
if (-not (Test-Path $profilePath)) {
    New-Item -ItemType File -Force -Path $profilePath | Out-Null
}

$sourceLine = "`$env:JVM_DIR = `"$jvmDir`"; if (Test-Path `"`$env:JVM_DIR\jvm.ps1`") { . `"`$env:JVM_DIR\jvm.ps1`" }"

$content = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
if ($content -notlike "*jvm.ps1*") {
    Add-Content -Path $profilePath -Value "`n# >>> jvm initialization >>>`n$sourceLine`n# <<< jvm initialization <<<"
    Write-Host "✓ Added jvm loader to $profilePath" -ForegroundColor Green
} else {
    Write-Host "✓ jvm is already configured in $profilePath" -ForegroundColor Green
}

Write-Host "`n🎉 jvm was successfully installed!" -ForegroundColor Green
Write-Host "To start using it now, run:`n  . `$PROFILE`n" -ForegroundColor Yellow
