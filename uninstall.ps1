# jvm Windows PowerShell Uninstaller
# Usage: irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.ps1 | iex

$ErrorActionPreference = "Stop"

$jvmDir = if ($env:JVM_DIR) { $env:JVM_DIR } else { "$HOME\.jvm" }

Write-Host "⚠️  This will uninstall jvm from your system." -ForegroundColor Yellow
$confirm = Read-Host "Are you sure you want to uninstall jvm? [y/N]"

if ($confirm -notmatch "^[yY]([eE][sS])?$") {
    Write-Host "Uninstallation cancelled." -ForegroundColor Cyan
    exit 0
}

# Nested prompt: Ask whether to keep or remove downloaded JDKs
$keepJdks = $true
$versionsDir = "$jvmDir\versions"
if ((Test-Path $versionsDir) -and (Get-ChildItem $versionsDir)) {
    Write-Host "`nFound installed JDK versions in $versionsDir." -ForegroundColor Yellow
    $jdkChoice = Read-Host "Do you want to keep your installed JDKs? [Y/n]"
    if ($jdkChoice -match "^[nN]([oO])?$") {
        $keepJdks = $false
    }
}

Write-Host "`nUninstalling jvm..." -ForegroundColor Cyan

# 1. Clean PowerShell profile
$profilePath = $PROFILE
if (Test-Path $profilePath) {
    $content = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($content -and ($content -like "*jvm*")) {
        Write-Host "Removing jvm configuration from $profilePath..." -ForegroundColor Cyan
        $newContent = $content -replace "(?ms)# >>> jvm initialization >>>.*?# <<< jvm initialization <<<", ""
        $newContent = $newContent -replace ".*jvm\.ps1.*", ""
        Set-Content -Path $profilePath -Value $newContent.Trim()
        Write-Host "✓ Cleaned $profilePath" -ForegroundColor Green
    }
}

# 2. Remove jvm directory or clean files
if (Test-Path $jvmDir) {
    if ($keepJdks) {
        Write-Host "Keeping installed JDKs in $versionsDir." -ForegroundColor Cyan
        Remove-Item -Recurse -Force -ErrorAction SilentlyContinue "$jvmDir\bin"
        Remove-Item -Recurse -Force -ErrorAction SilentlyContinue "$jvmDir\downloads"
        Remove-Item -Force -ErrorAction SilentlyContinue "$jvmDir\current"
        Remove-Item -Force -ErrorAction SilentlyContinue "$jvmDir\jvm.ps1"
        Remove-Item -Force -ErrorAction SilentlyContinue "$jvmDir\config.json"
        Write-Host "✓ Removed jvm binary, scripts, and runtime files" -ForegroundColor Green
    } else {
        Write-Host "Removing all jvm files and JDK versions from $jvmDir..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force -Path $jvmDir
        Write-Host "✓ Removed $jvmDir" -ForegroundColor Green
    }
}

Write-Host "`n🎉 jvm was successfully uninstalled!" -ForegroundColor Green
Write-Host "Please restart your PowerShell session to finish.`n"
