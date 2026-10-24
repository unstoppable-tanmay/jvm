# jvm (Java Version Manager)

`jvm` is an ultra-fast, cross-platform Java Version Manager written in Go with an `nvm`-style architecture. It simplifies discovering, installing, switching, linking, and managing multiple JDKs without administrative privileges or messy manual environment variable setups.

---

## ⚡ Architecture (Identical to `nvm`)

Just like `nvm`:
- **`install.sh` / `install.ps1`**: The installation script copies the binary and places the loader in your shell profile.
- **`jvm.sh` (or `jvm.ps1` / `jvm.fish`)**: Sourced once in your `.zshrc` / `.bashrc` / `$PROFILE`, providing the thin `jvm` wrapper function that instantly updates `JAVA_HOME` and `PATH` on `jvm use`.
- **`jvm` binary**: High-speed compiled Go binary that handles all downloading, checksum verification, archive extraction, and version resolution.

---

## 🚀 Quick Install

### macOS / Linux (POSIX)
```bash
# From git repository:
./install.sh

# Or via curl:
# curl -fsSL https://raw.githubusercontent.com/tanmay/jvm/main/install.sh | bash
```

Reload your shell:
```bash
source ~/.zshrc   # or source ~/.bashrc
```

### Windows (PowerShell)
```powershell
.\install.ps1
# Or: irm https://raw.githubusercontent.com/tanmay/jvm/main/install.ps1 | iex
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

When you `cd` into that directory, `jvm.sh` automatically switches your `JAVA_HOME` and `PATH` to that version!
