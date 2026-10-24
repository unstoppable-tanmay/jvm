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
            Invoke-Expression $envScript
        }
    } else {
        & "jvm.exe" @args
    }
}

$initEnv = & "jvm.exe" env --shell powershell 2>$null
if ($initEnv) { Invoke-Expression $initEnv }
