#!/usr/bin/env bash
set -e

# jvm installation script (curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash)

JVM_DIR="${JVM_DIR:-$HOME/.jvm}"
GITHUB_REPO="unstoppable-tanmay/jvm"

echo "Installing jvm to $JVM_DIR..."

mkdir -p "$JVM_DIR/bin"
mkdir -p "$JVM_DIR/versions"
mkdir -p "$JVM_DIR/downloads"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd)"

# If run locally from repository root
if [ -f "$SCRIPT_DIR/bin/jvm" ]; then
    cp "$SCRIPT_DIR/bin/jvm" "$JVM_DIR/bin/jvm"
    chmod +x "$JVM_DIR/bin/jvm"
    cp "$SCRIPT_DIR/jvm.sh" "$JVM_DIR/jvm.sh"
    [ -f "$SCRIPT_DIR/jvm.fish" ] && cp "$SCRIPT_DIR/jvm.fish" "$JVM_DIR/jvm.fish"
elif [ -f "./bin/jvm" ]; then
    cp "./bin/jvm" "$JVM_DIR/bin/jvm"
    chmod +x "$JVM_DIR/bin/jvm"
    cp "./jvm.sh" "$JVM_DIR/jvm.sh"
    [ -f "./jvm.fish" ] && cp "./jvm.fish" "$JVM_DIR/jvm.fish"
else
    # Remote release installer: detect OS and Arch
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
    esac

    echo "Fetching latest release for $OS-$ARCH..."
    LATEST_TAG=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_TAG" ]; then
        echo "Could not resolve latest release tag. Falling back to downloading jvm.sh..."
    else
        ARCHIVE_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_TAG/jvm_${LATEST_TAG#v}_${OS}_${ARCH}.tar.gz"
        echo "Downloading $ARCHIVE_URL..."
        TEMP_TAR="/tmp/jvm_${LATEST_TAG}.tar.gz"
        curl -fsSL "$ARCHIVE_URL" -o "$TEMP_TAR"
        tar -xzf "$TEMP_TAR" -C "$JVM_DIR/bin" jvm 2>/dev/null || tar -xzf "$TEMP_TAR" -C "$JVM_DIR/bin"
        chmod +x "$JVM_DIR/bin/jvm"
        rm -f "$TEMP_TAR"
    fi

    echo "Downloading shell integration scripts..."
    curl -fsSL "https://raw.githubusercontent.com/$GITHUB_REPO/main/jvm.sh" -o "$JVM_DIR/jvm.sh"
    curl -fsSL "https://raw.githubusercontent.com/$GITHUB_REPO/main/jvm.fish" -o "$JVM_DIR/jvm.fish" 2>/dev/null || true
fi

# Detect user profile
SHELL_NAME="$(basename "$SHELL" 2>/dev/null || echo "bash")"
PROFILE=""

if [ "$SHELL_NAME" = "zsh" ]; then
    PROFILE="$HOME/.zshrc"
elif [ "$SHELL_NAME" = "bash" ]; then
    if [ -f "$HOME/.bash_profile" ]; then
        PROFILE="$HOME/.bash_profile"
    else
        PROFILE="$HOME/.bashrc"
    fi
elif [ "$SHELL_NAME" = "fish" ]; then
    PROFILE="$HOME/.config/fish/config.fish"
fi

SOURCE_STR="[ -s \"$JVM_DIR/jvm.sh\" ] && \\. \"$JVM_DIR/jvm.sh\" # This loads jvm"

if [ -n "$PROFILE" ]; then
    mkdir -p "$(dirname "$PROFILE")"
    touch "$PROFILE"
    if ! grep -q "jvm.sh" "$PROFILE"; then
        echo "" >> "$PROFILE"
        echo "# >>> jvm initialization >>>" >> "$PROFILE"
        echo "export JVM_DIR=\"$JVM_DIR\"" >> "$PROFILE"
        echo "$SOURCE_STR" >> "$PROFILE"
        echo "# <<< jvm initialization <<<" >> "$PROFILE"
        echo "✓ Added jvm loader to $PROFILE"
    else
        echo "✓ jvm is already configured in $PROFILE"
    fi
fi

echo ""
echo "🎉 jvm was successfully installed!"
echo "To start using it now, run:"
[ -n "$PROFILE" ] && echo "  source $PROFILE" || echo "  source $JVM_DIR/jvm.sh"
echo ""
