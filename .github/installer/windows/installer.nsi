!include "MUI2.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"

Name "jvm (Java Version Manager)"
OutFile "jvm.exe"
InstallDir "$PROFILE\.jvm"
RequestExecutionLevel user

VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "jvm"
VIAddVersionKey "Comments" "Java version manager for Windows"
VIAddVersionKey "CompanyName" "unstoppable-tanmay"
VIAddVersionKey "FileDescription" "jvm Setup (Windows x64)"
VIAddVersionKey "FileVersion" "1.0.0.0"
VIAddVersionKey "LegalCopyright" "MIT License"

Var Dialog
Var JavaDirInput
Var JavaDirText
Var LinkDirInput
Var LinkDirText

!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
Page custom CustomOptionsPage CustomOptionsPageLeave
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Function CustomOptionsPage
    nsDialogs::Create 1018
    Pop $Dialog
    ${If} $Dialog == error
        Abort
    ${EndIf}

    ${NSD_CreateLabel} 0 0 100% 24u "Custom JDK Storage Location (Optional):$\nLeave empty to use default ($INSTDIR\versions) or specify an external drive:"
    Pop $0

    ${NSD_CreateText} 0 26u 100% 12u "$JavaDirText"
    Pop $JavaDirInput

    ${NSD_CreateLabel} 0 46u 100% 24u "Link Existing Java Installation (Optional):$\nIf you already have a JDK on your machine, enter its path (e.g. C:\Program Files\Java\jdk-17):"
    Pop $0

    ${NSD_CreateText} 0 72u 100% 12u "$LinkDirText"
    Pop $LinkDirInput

    nsDialogs::Show
FunctionEnd

Function CustomOptionsPageLeave
    ${NSD_GetText} $JavaDirInput $JavaDirText
    ${NSD_GetText} $LinkDirInput $LinkDirText
FunctionEnd

Section "Install"
    SetOutPath "$INSTDIR\bin"
    !ifdef BIN_PATH
        File "/oname=jvm.exe" "${BIN_PATH}"
    !else
        File "dist\jvm_windows_amd64_v1\jvm.exe"
    !endif
    File "jvm.cmd"

    SetOutPath "$INSTDIR"
    File "jvm.ps1"
    File "README.md"
    File "LICENSE"

    CreateDirectory "$INSTDIR\versions"
    CreateDirectory "$INSTDIR\downloads"

    ; If custom Java storage directory was specified, create it and write config.json
    ${If} $JavaDirText != ""
        CreateDirectory "$JavaDirText"
        FileOpen $0 "$INSTDIR\config.json" w
        FileWrite $0 '{"storage_dir": "$JavaDirText"}'
        FileClose $0
    ${EndIf}

    ; Add $INSTDIR\bin to User PATH permanently via PowerShell (reliable & built-in)
    ExecWait 'powershell -NoProfile -ExecutionPolicy Bypass -Command "$p = [Environment]::GetEnvironmentVariable(\"Path\", \"User\"); if ($p -notlike \"*$INSTDIR\bin*\") { [Environment]::SetEnvironmentVariable(\"Path\", \"$INSTDIR\bin;\" + $p, \"User\") }"'

    ; If user specified an existing JDK to link, run jvm link
    ${If} $LinkDirText != ""
        ExecWait '"$INSTDIR\bin\jvm.exe" link "$LinkDirText"'
    ${EndIf}

    ; Generate standard shims
    ExecWait '"$INSTDIR\bin\jvm.exe" use default'

    ; Write uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
    ExecWait 'powershell -NoProfile -ExecutionPolicy Bypass -Command "$p = [Environment]::GetEnvironmentVariable(\"Path\", \"User\"); $new = ($p.Split(\";\") | Where-Object { $_ -ne \"$INSTDIR\bin\" }) -join \";\"; [Environment]::SetEnvironmentVariable(\"Path\", $new, \"User\")"'

    Delete "$INSTDIR\bin\jvm.exe"
    Delete "$INSTDIR\bin\jvm.cmd"
    Delete "$INSTDIR\bin\java.cmd"
    Delete "$INSTDIR\bin\javac.cmd"
    Delete "$INSTDIR\bin\*.cmd"
    Delete "$INSTDIR\bin\*.ps1"
    RMDir "$INSTDIR\bin"
    Delete "$INSTDIR\jvm.ps1"
    Delete "$INSTDIR\README.md"
    Delete "$INSTDIR\LICENSE"
    Delete "$INSTDIR\config.json"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"
SectionEnd
