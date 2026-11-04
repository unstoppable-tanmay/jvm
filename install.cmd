@echo off
REM jvm Windows Batch Installer
REM Can be run from Command Prompt or downloaded and double-clicked.
setlocal enabledelayedexpansion

echo Installing jvm...

if "%JVM_DIR%"=="" (
    set "JVM_DIR=%USERPROFILE%\.jvm"
)
set "BIN_DIR=%JVM_DIR%\bin"

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%JVM_DIR%\versions" mkdir "%JVM_DIR%\versions"

powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "$repo = 'unstoppable-tanmay/jvm'; " ^
    "$release = Invoke-RestMethod -Uri \"https://api.github.com/repos/$repo/releases/latest\" -Headers @{ 'User-Agent' = 'jvm-installer' }; " ^
    "$tag = $release.tag_name; $version = $tag.TrimStart('v'); " ^
    "$url = \"https://github.com/$repo/releases/download/$tag/jvm_${version}_windows_amd64.zip\"; " ^
    "$zip = \"$env:TEMP\jvm.zip\"; " ^
    "Invoke-WebRequest -Uri $url -OutFile $zip; " ^
    "Expand-Archive -Path $zip -DestinationPath '%BIN_DIR%' -Force; " ^
    "Remove-Item $zip -Force; " ^
    "Invoke-WebRequest -Uri \"https://raw.githubusercontent.com/$repo/main/jvm.ps1\" -OutFile '%JVM_DIR%\jvm.ps1'; " ^
    "Invoke-WebRequest -Uri \"https://raw.githubusercontent.com/$repo/main/jvm.cmd\" -OutFile '%BIN_DIR%\jvm.cmd'; " ^
    "$userPath = [Environment]::GetEnvironmentVariable('Path', 'User'); " ^
    "if ($userPath -notlike '*%BIN_DIR%*') { [Environment]::SetEnvironmentVariable('Path', '%BIN_DIR%;' + $userPath, 'User') }"

echo.
echo jvm successfully installed to %JVM_DIR%
echo Restart your Command Prompt or PowerShell to start using jvm!

