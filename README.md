# jvm (Java Version Manager)

`jvm` is an ultra-fast, cross-platform Java Version Manager written in Go with an `nvm`-style architecture. It simplifies discovering, installing, switching, linking, and managing multiple JDKs without administrative privileges or messy manual environment variable setups.

---

## 💡 Why `jvm`?

- **⚡ Blazing Fast Shell Startup**: Sourcing heavy shell managers can noticeably delay opening new terminal tabs. `jvm` uses a high-performance compiled Go core paired with a minimal ~15-line shell wrapper, keeping shell startup times under **3ms**.
- **📦 Zero External Dependencies**: A single self-contained, statically compiled binary. Does not depend on system `curl`, `tar`, `unzip`, `jq`, Python, or Node. Works identically across minimal Docker images, Linux, macOS, and Windows.
- **🪟 True Cross-Platform & Native Windows Support**: First-class support for macOS (Apple Silicon `arm64` & Intel `x64`), Linux (`amd64` & `arm64`), and Windows with native PowerShell integration—no WSL or Cygwin required.
- **🔗 Seamless External JDK Linking**: Connect JDKs already present on your machine (e.g. from Homebrew, Apple native installers, or Android Studio) without redownloading 300MB archives. `jvm link` automatically parses vendor metadata and links them safely.
- **🍎 Smart macOS Bundle Awareness**: Automatically navigates the differences between standard Linux root layouts and macOS `.jdk/Contents/Home` structures to guarantee correct `JAVA_HOME` resolution across all tools (IntelliJ, Gradle, Maven).
- **💾 Configurable Storage & Directory Auto-Scanning**: Store JDKs on external drives, secondary partitions, or custom enterprise volumes with `jvm set-dir`, automatically discovering and registering all existing installations in that path.
- **🌐 Install from Anywhere**: In addition to querying all major upstream distributions (Temurin, Corretto, Zulu, GraalVM), download and verify custom releases, early-access builds (EA), or private enterprise mirrors directly via `--from <url>`.
- **📁 Frictionless Project-Level Switching**: Automatically switch `JAVA_HOME` and `$PATH` upon entering a directory containing `.java-version` or `.jvmrc`.

---

## ⚡ Architecture

- **`install.sh` / `install.ps1`**: Installs the binary and places the loader in your shell profile.
- **`jvm.sh` (or `jvm.ps1` / `jvm.fish`)**: Lightweight shell wrapper that instantly updates `JAVA_HOME` and `PATH` on `jvm use`.
- **`jvm` binary**: Compiled Go binary handling discovery, downloads, SHA256 verification, archive extraction, and version resolution.

---

## 🚀 Quick Install

### macOS / Linux (POSIX)
```bash
# Via curl:
curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash

# Or from local repository:
./install.sh
```

Reload your shell:
```bash
source ~/.zshrc   # or source ~/.bashrc
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1 | iex

# Or from local repository:
.\install.ps1
```

---

## 🛠️ Usage

### 1. Switch Versions (Instant & Seamless)
```bash
# Switches active version in current terminal session immediately
jvm use 21
jvm use corretto@17
```

### 2. Install a JDK
```bash
# Install Eclipse Temurin 21 (default vendor)
jvm install 21

# Install Amazon Corretto 17
jvm install corretto@17

# Install from custom direct URL or early-access build
jvm install --from "https://github.com/.../jdk.tar.gz" custom-21
```

### 3. List Installed and Remote Versions
```bash
# List local installed and linked JDKs
jvm ls

# List remote available versions
jvm ls-remote --vendor temurin --version 21

# Show active version and path
jvm current
```

### 4. Link an Existing JDK
```bash
# Link by path (auto-detects version & vendor from release metadata)
jvm link /Library/Java/JavaVirtualMachines/temurin-17.jdk

# Or link with a custom alias
jvm link /opt/homebrew/opt/openjdk@11 homebrew-11
```

### 5. Set Custom Directory & Auto-Scan
```bash
# Point jvm to store JDKs in a custom directory (e.g. an external drive)
# Also automatically scans and registers any existing JDKs found inside that path!
jvm set-dir /Volumes/ExternalSSD/jdks
```

### 6. Uninstall or Unlink
```bash
# Uninstalls downloaded JDKs, or cleanly unlinks external ones (without deleting files)
jvm uninstall 17
```

---

## 📁 Per-Project Version Switching (`.java-version`)

To lock a project to a specific Java version, create a `.java-version` or `.jvmrc` file in your repository root:

```bash
echo "21" > .java-version
```

When you `cd` into that directory, `jvm` automatically switches your `JAVA_HOME` and `PATH` to that version!
