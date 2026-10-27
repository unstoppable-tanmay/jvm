# jvm (Java Version Manager)

[![GitHub Sponsors](https://img.shields.io/badge/Sponsor-%E2%99%A5-ea4aaa?style=flat&logo=githubsponsors&logoColor=white)](https://github.com/sponsors/unstoppable-tanmay)
[![Release](https://img.shields.io/github/v/release/unstoppable-tanmay/jvm?style=flat&color=18181b)](https://github.com/unstoppable-tanmay/jvm/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Website](https://img.shields.io/badge/docs-website-zinc?style=flat)](https://github.com/unstoppable-tanmay/jvm)

> **Java version management, made simple.**  
> Stop wrestling with broken `JAVA_HOME` paths, slow shell startups, and manual downloads. Switch between Java 8, 17, 21, and beyond in milliseconds.

---

## ⚡ Highlights

- **🚀 Sub-3ms Startup**: Sourcing heavy shell managers can noticeably delay opening new terminal tabs. `jvm` uses a compiled Go core with a ~15-line shell wrapper to keep shell startup under 3ms.
- **🔗 Keep Your Existing JDKs**: Connect JDKs already present on your machine (Homebrew, Apple native installers, Android Studio) without redownloading 300MB archives.
- **📁 Frictionless Project Switching**: Commit a `.java-version` or `.jvmrc` file to your repository. When you or your teammates `cd` in, the right Java runtime activates automatically.
- **🌐 Install from Anywhere**: Query major upstream distributions (Temurin, Corretto, Zulu, GraalVM), or download custom internal enterprise mirrors and early-access builds with `--from <url>` and SHA256 verification.
- **🪟 True Cross-Platform**: First-class support for macOS (Apple Silicon `arm64` & Intel `x64`), Linux (`amd64` & `arm64`), and Windows (native PowerShell & `cmd.exe`).
- **💾 Flexible Storage**: Store JDKs on external drives, secondary partitions, or custom enterprise volumes with `jvm set-dir` and safe migration/linking prompts.

---

## 🚀 Quick Install

### macOS / Linux (curl or wget)
```bash
# Using curl:
curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash

# Or using wget:
wget -qO- https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash
```

### Homebrew (macOS / Linux)
```bash
brew tap unstoppable-tanmay/jvm
brew install jvm
```
Add the loader to your shell profile (`~/.zshrc` or `~/.bashrc`):
```bash
[ -s "$(brew --prefix jvm)/jvm.sh" ] && \. "$(brew --prefix jvm)/jvm.sh"
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1 | iex
```

### Direct Binary Download
Pre-compiled standalone archives for every OS and architecture are available on the [**Releases Page**](https://github.com/unstoppable-tanmay/jvm/releases).

---

## 🛠️ Command Reference

| Command | Shorthand | Description |
|---|---|---|
| `jvm use <version>` | — | Switches active JDK in current shell tab immediately (prompts selector if multiple match) |
| `jvm install <version>` | `jvm i` | Downloads and verifies JDK from upstream vendors (Temurin, Corretto, Zulu, GraalVM) |
| `jvm ls` | — | Lists all locally installed and externally linked JDK versions |
| `jvm ls-remote` | — | Queries available upstream releases by vendor and version |
| `jvm link <path> [name]` | — | Links an existing external JDK into `jvm` without redownloading |
| `jvm unlink <name>` | — | Safely unlinks an external JDK without touching original disk files |
| `jvm set-dir <path>` | — | Points storage to another folder (e.g. external SSD) with migration prompts |
| `jvm current` | — | Displays active `jvm` Java version and system runtime path |
| `jvm uninstall <version>` | `jvm ui` | Removes a downloaded JDK from `jvm` storage |

### Common Examples

#### Switch Java Runtime
```bash
# Switch to version 21
jvm use 21

# Switch to specific vendor distribution
jvm use corretto@17

# If multiple versions match, an intuitive selector appears:
# Multiple matching JDKs found for switching active JDK:
#   [1] temurin-21+35   (installed)
#   [2] homebrew-21     (linked -> /opt/homebrew)
# Select an option [1-2]: 1
```

#### Install OpenJDK
```bash
# Install Eclipse Temurin 21 (default vendor)
jvm install 21
# Shorthand alias
jvm i 21

# Install Amazon Corretto 17
jvm i corretto@17

# Install from custom direct URL with SHA256 integrity check
jvm install --from "https://github.com/.../jdk.tar.gz" --sha256 "<hash>" custom-21
```

#### Link Existing JDKs
```bash
# Link macOS native installer bundle (auto-detects release metadata)
jvm link /Library/Java/JavaVirtualMachines/temurin-17.jdk

# Link Homebrew JDK with custom alias
jvm link /opt/homebrew/opt/openjdk@11 homebrew-11

# Unlink safely (original files remain untouched)
jvm unlink homebrew-11
```

#### Custom Storage Directory
```bash
# Change JDK storage directory
jvm set-dir /Volumes/ExternalSSD/jdks

# Prompts:
# Existing JDKs found in ~/.jvm/versions:
#   [1] Move (migrate all existing JDKs to new directory)
#   [2] Link (symlink existing JDKs into new directory)
#   [3] Leave (keep them where they are without linking)
```

---

## 📁 Automatic Project Switching

Lock any project to a specific Java version by adding a `.java-version` or `.jvmrc` file in your repository:

```bash
echo "21" > .java-version
```

When you `cd` into any repository or subdirectory containing `.java-version`, `jvm` walks up the tree, matches the version, and sets `JAVA_HOME` and `PATH` automatically.

---

## ✅ Features & Roadmap

### User-Facing Capabilities
- [x] Instant version switching in your current shell session without restarting terminal
- [x] Auto-switch Java runtime when navigating into projects with `.java-version` or `.jvmrc`
- [x] Link existing JDKs from Homebrew, Android Studio, and macOS installers without redownloading
- [x] Clean unlinking to remove links without touching original external files on disk
- [x] Interactive terminal selector when multiple installed or linked versions match a target
- [x] Install directly from upstream vendors (Temurin, Corretto, Zulu, GraalVM) via Disco API
- [x] Install from custom direct URLs with optional SHA256 integrity verification (`--from`)
- [x] Configurable storage location (e.g. external SSDs) with safe migration prompts (`jvm set-dir`)
- [x] Cross-platform support across macOS (Apple Silicon & Intel), Linux, and Windows (PowerShell & `cmd.exe`)
- [x] Clean interactive uninstallation with prompt to preserve or delete downloaded JDKs

### Upcoming Roadmap
- [ ] Automatic Homebrew formula updates on new GitHub release tags
- [ ] Built-in proxy and enterprise mirror registry support in user settings
- [ ] Direct migration wizard to import existing SDKMAN candidate directories

---

## 🗑️ Uninstalling `jvm`

To completely remove `jvm` and clean up your shell profiles:

### macOS / Linux (POSIX)
```bash
curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.ps1 | iex
```

> **Note**: The uninstaller prompts for confirmation and asks whether you want to **keep your installed JDK versions** or delete everything completely.

---

## ☕ Support & Sponsor

If you find `jvm` useful, consider sponsoring to support continuous maintenance and new features!

<a href="https://github.com/sponsors/unstoppable-tanmay">
  <img src="https://img.shields.io/badge/Sponsor_on_GitHub-%E2%99%A5-ea4aaa?style=for-the-badge&logo=githubsponsors&logoColor=white" alt="Sponsor unstoppable-tanmay" />
</a>

---

## 📄 License

MIT License © [Tanmay Kumar Panda](https://tanmaykumarpanda.in)
