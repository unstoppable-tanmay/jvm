# jvm PowerShell integration
if (-not $env:JVM_DIR) {
    $env:JVM_DIR = "$HOME\.jvm"
}

$jvmBinDir = "$env:JVM_DIR\bin"
if ((Test-Path $jvmBinDir) -and ($env:PATH -notlike "*$jvmBinDir*")) {
    $env:PATH = "$jvmBinDir;" + $env:PATH
}

function jvm {
    if ($args.Count -gt 0 -and $args[0] -eq "use") {
        & "jvm.exe" @args
        if ($LASTEXITCODE -eq 0) {
            $envScript = & "jvm.exe" env --shell powershell
            if ($envScript) { Invoke-Expression $envScript }
        }
    } else {
        & "jvm.exe" @args
    }
}

# Auto-switch Java version when changing directory if .java-version or .jvmrc is present
function _jvm_auto_switch {
    $checkDir = (Get-Location).Path
    $found = $false
    while ($checkDir -and ($checkDir -ne [System.IO.Path]::GetPathRoot($checkDir))) {
        if ((Test-Path (Join-Path $checkDir ".java-version")) -or (Test-Path (Join-Path $checkDir ".jvmrc"))) {
            $found = $true
            break
        }
        $checkDir = Split-Path $checkDir -Parent
    }
    if ($found) {
        $autoEnv = & "jvm.exe" env --shell powershell 2>$null
        if ($autoEnv) { Invoke-Expression $autoEnv }
    }
}

# Hook PowerShell prompt to run auto-switch check on every directory change
if (Test-Path Function:\prompt) {
    $oldPrompt = (Get-Command prompt).Definition
    function prompt {
        _jvm_auto_switch
        & ([ScriptBlock]::Create($oldPrompt))
    }
} else {
    function prompt {
        _jvm_auto_switch
        "PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) "
    }
}

$initEnv = & "jvm.exe" env --shell powershell 2>$null
if ($initEnv) { Invoke-Expression $initEnv }
