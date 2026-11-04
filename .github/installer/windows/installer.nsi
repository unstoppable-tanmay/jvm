!include "MUI2.nsh"
!include "EnvVarUpdate.nsh"

Name "jvm (Java Version Manager)"
OutFile "jvm-setup.exe"
InstallDir "$PROFILE\.jvm"
RequestExecutionLevel user

!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

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

    ; Add $INSTDIR\bin to User PATH permanently
    ${EnvVarUpdate} $0 "PATH" "A" "HKCU" "$INSTDIR\bin"

    ; Write uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
    ; Remove from User PATH
    ${un.EnvVarUpdate} $0 "PATH" "R" "HKCU" "$INSTDIR\bin"

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
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"
SectionEnd
