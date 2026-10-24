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
} elseif (Test-Path ".\bin\jvm.exe") {
    Copy-Item ".\bin\jvm.exe" -Destination "$binDir\jvm.exe" -Force
    Copy-Item ".\jvm.ps1" -Destination "$jvmDir\jvm.ps1" -Force
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
