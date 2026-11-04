# AGENT ARCHITECTURE & CONTEXT: `jvm` (Java Version Manager)

## Project Overview
`jvm` is a cross-platform Java Version Manager written in Go with lightweight shell integration scripts for macOS, Linux, and Windows. It provides fast runtime switching, upstream JDK downloads (via Foojay Disco API and Adoptium), directory-based version switching (`.java-version`, `.jvmrc`), and external JDK linking.

---

## Directory Structure & Key Files
- `cmd/jvm/main.go`: Entry point for the `jvm` CLI binary.
- `internal/cli/`:
  - `root.go`: Root Cobra command and subcommands registration.
  - `env.go`: Generates shell export commands (`export JAVA_HOME=...; export PATH=...`).
  - `use.go`: Handles `jvm use <version>` switching active JDK via symlink/junction/shims.
  - `install.go`: Installs JDKs from upstream or custom direct URLs (`--from`), with automatic fallback for vendor-platform mismatches.
  - `list_remote.go`: Queries remote JDK versions (`jvm ls-remote`) with `--search` multi-filtering and version matching.
  - `list.go`: Lists installed and linked JDKs (`jvm ls`).
  - `link.go` / `unlink.go`: Links external JDKs without copying files.
  - `current.go`: Displays active version and runtime path.
- `internal/disco/client.go`: Foojay Disco v3.0 API client for querying packages and download URLs. Handles `version_number` vs `jdk_version` for Java 8.
- `internal/config/`:
  - `config.go`: Resolves directories (`~/.jvm`, `bin`, `versions`, `downloads`, `current`).
  - `link.go`: Handles cross-platform current version linking (`Symlink`, directory junctions `mklink /J`, and `current.path` fallback for Windows without Developer Mode / Admin).
  - `shims.go`: Automatic shims in `~/.jvm/bin` for `java`, `javac`, `jar`, etc., for non-interactive shells and background jobs.
- `internal/platform/platform.go`: OS and Architecture detection (`macos`, `linux`, `windows`; `amd64`, `arm64`).
- `jvm.sh`: Shell wrapper for bash/zsh (handles function wrapping, `chpwd` / `PROMPT_COMMAND` auto-switch hooks).
- `jvm.ps1`: PowerShell wrapper module with auto-switch prompt hook on `cd`.
- `jvm.cmd`: Batch script wrapper for Windows CMD (`cmd.exe`).
- `jvm.fish`: Fish shell wrapper.
- `install.sh`: Bash installer for Unix-like systems.
- `install.ps1`: PowerShell installer for Windows (updates User `PATH` permanently in Registry).

---

## Solutions Implemented

### 1. Non-Interactive Terminals
- Implemented `EnsureShims` in [internal/config/shims.go](file:///Users/tanmay/workspace/projects/jvm/internal/config/shims.go).
- Generates shims in `~/.jvm/bin/` (`java`, `javac`, `jar`, `keytool`, `jshell`, etc.) that automatically resolve the active Java installation from `current` or `current.path`.
- Non-interactive terminals, CI/CD runners, cron jobs, and subshells can execute `java` directly as long as `~/.jvm/bin` is in PATH.

### 2. Windows Native Support & Auto-switching
- **PowerShell Auto-switching**: Hooked `prompt` function in [jvm.ps1](file:///Users/tanmay/workspace/projects/jvm/jvm.ps1) to detect `.java-version` / `.jvmrc` up the directory tree and apply version switches dynamically upon `cd`.
- **PowerShell / CMD PATH persistence**: Updated [install.ps1](file:///Users/tanmay/workspace/projects/jvm/install.ps1) to add `$binDir` permanently to User `PATH` via `[EnvironmentVariableTarget]::User`, so CMD and PowerShell recognize both `jvm` and `java` without requiring WSL.
- **Windows Junctions & Fallbacks**: In [internal/config/link.go](file:///Users/tanmay/workspace/projects/jvm/internal/config/link.go), added Directory Junction (`mklink /J`) and `current.path` fallback if Windows Developer Mode / Admin rights are not active.
- **Installer Clarification**: Updated docs to explain `irm` (`Invoke-RestMethod`) and `iex` (`Invoke-Expression`), plus added a zero-dependency bypass command:
  `powershell -ExecutionPolicy Bypass -Command "Invoke-Expression (Invoke-RestMethod 'https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1')"`

### 3. Java 8 Installation & `ls-remote` Search
- In [internal/disco/client.go](file:///Users/tanmay/workspace/projects/jvm/internal/disco/client.go), fixed Foojay parameter handling: queries with pure integers (e.g. `8`, `17`, `21`) use `jdk_version=<num>` directly, resolving all builds for Java 8 across all supported vendors.
- In [internal/cli/list_remote.go](file:///Users/tanmay/workspace/projects/jvm/internal/cli/list_remote.go), removed the hardcoded default vendor (`temurin`). `jvm ls-remote` now performs a global multi-vendor search across all distributions by default unless `-v <vendor>` is explicitly specified.
- In [internal/cli/install.go](file:///Users/tanmay/workspace/projects/jvm/internal/cli/install.go), added smart distribution fallback: if Temurin has no build for that architecture (e.g., Temurin has no native `macos/arm64` Java 8 build), it automatically checks Corretto and Zulu and selects the valid release.
- In [internal/cli/list_remote.go](file:///Users/tanmay/workspace/projects/jvm/internal/cli/list_remote.go), added `--search` (`-s`) supporting multi-keyword filtering (comma or space separated, e.g. `--search "corretto,fx"` or `"8,lts"`).

### 4. Release Packaging (.dmg & .exe)
- **Roadmap / In Progress**:
  - Windows: Setup `.exe` / MSI installer (via NSIS or WiX in GitHub Actions) that installs to `%USERPROFILE%\.jvm`, adds `%USERPROFILE%\.jvm\bin` to User `PATH`, and sets up shims.
  - macOS: `.dmg` drag-and-drop or `.pkg` installer that places `jvm` and shell integration files directly into system/user path.



