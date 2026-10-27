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
- **`jvm.sh` (or `jvm.ps1` / `jvm.fish`)**: Lightweight shell wrapper that instantly updates `JAVA_HOME` and `PATH` on `jvm use` and hooks into shell directory transitions (`cd`).
- **`jvm` binary**: Compiled Go binary handling discovery, downloads, SHA256 verification, archive extraction, and version resolution.

---

## 🚀 Installation Methods

### 1. Standalone Script (Recommended)

#### macOS & Linux (curl or wget)
```bash
# Using curl:
curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash

# Or using wget:
wget -qO- https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash
```

#### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1 | iex
```

---

### 2. Homebrew (macOS / Linux)
```bash
brew tap unstoppable-tanmay/jvm
brew install jvm
```

Add the loader to your shell profile (`~/.zshrc` or `~/.bashrc`):
```bash
[ -s "$(brew --prefix jvm)/jvm.sh" ] && \. "$(brew --prefix jvm)/jvm.sh"
```

---

### 3. Via `go install`
If you already have Go installed on your machine:
```bash
go install github.com/unstoppable-tanmay/jvm/cmd/jvm@latest
```

---

### 4. Direct Binary Download (GitHub Releases)
Pre-compiled standalone archives for every OS and architecture are available directly on the [**GitHub Releases Page**](https://github.com/unstoppable-tanmay/jvm/releases):

| OS | Architecture | Format | Archive Link |
|---|---|---|---|
| **macOS** | Apple Silicon (`arm64`) | `.tar.gz` | [Download](https://github.com/unstoppable-tanmay/jvm/releases/latest) |
| **macOS** | Intel (`x86_64`) | `.tar.gz` | [Download](https://github.com/unstoppable-tanmay/jvm/releases/latest) |
| **Linux** | 64-bit (`amd64`) | `.tar.gz` | [Download](https://github.com/unstoppable-tanmay/jvm/releases/latest) |
| **Linux** | ARM64 (`aarch64`) | `.tar.gz` | [Download](https://github.com/unstoppable-tanmay/jvm/releases/latest) |
| **Windows**| 64-bit (`x64`) | `.zip` | [Download](https://github.com/unstoppable-tanmay/jvm/releases/latest) |

#### Manual Setup for Downloaded Binaries:
1. Extract the downloaded archive.
2. Move the `jvm` binary to a folder in your `$PATH` (e.g. `~/.jvm/bin/` or `/usr/local/bin/`).
3. Source `jvm.sh` in your `~/.zshrc` or `~/.bashrc`:
   ```bash
   export JVM_DIR="$HOME/.jvm"
   [ -s "$JVM_DIR/jvm.sh" ] && \. "$JVM_DIR/jvm.sh"
   ```

---

### 5. Build from Source
```bash
git clone https://github.com/unstoppable-tanmay/jvm.git
cd jvm
./install.sh
```

---

## 🛠️ Usage

### 1. Switch Versions (Instant & Seamless)
```bash
# Switches active version in current terminal session immediately
jvm use 21
jvm use corretto@17
```
> If multiple JDKs match the target (e.g. an installed Temurin 17 and a linked Homebrew 17), `jvm` displays an interactive prompt asking you to select your preferred choice.

### 2. Install a JDK (alias: `i`)
```bash
# Install Eclipse Temurin 21 (default vendor)
jvm install 21
# Or using the shorthand alias:
jvm i 21

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
*(Linking internal `jvm` storage folders is disallowed to prevent cyclic loops).*

### 5. Unlink an External JDK
```bash
# Unlink an external JDK without touching the original files on disk
jvm unlink homebrew-11
```

### 6. Set Custom Directory & Auto-Scan
```bash
# Point jvm to store JDKs in a custom directory (e.g. an external drive)
# Also automatically scans and registers any existing JDKs found inside that path!
jvm set-dir /Volumes/ExternalSSD/jdks
```

### 7. Uninstall a Downloaded JDK (alias: `ui`)
```bash
# Uninstalls a downloaded JDK from jvm storage
jvm uninstall 17
# Or using the shorthand alias:
jvm ui 17
```

---

## 📁 Per-Project Version Switching (`.java-version` / `.jvmrc`)

To lock a project to a specific Java version, create a `.java-version` or `.jvmrc` file in your repository:

```bash
echo "21" > .java-version
```

### How Auto-Switching Works:
The `jvm.sh` / `jvm.fish` shell wrapper intercepts directory changes (`chpwd` hook in Zsh, `PROMPT_COMMAND` in Bash, `--on-variable PWD` in Fish). When you `cd` into any project containing `.java-version` (or any of its child subdirectories), `jvm` walks up the directory tree, matches the specified version against your installed/linked JDKs, and automatically points `JAVA_HOME` and `PATH` to it.

---

## 🗑️ Uninstalling `jvm`

To completely remove `jvm` and clean up your shell profiles:

### macOS / Linux (POSIX)
```bash
# Via curl:
curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.sh | bash

# Or from local repository:
./uninstall.sh
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.ps1 | iex

# Or from local repository:
.\uninstall.ps1
```

> **Note**: The uninstaller prompts for confirmation and asks whether you want to **keep your installed JDK versions** or delete everything completely.
