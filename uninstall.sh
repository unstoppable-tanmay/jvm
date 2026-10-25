#!/usr/bin/env bash
set -e

# jvm uninstallation script
# Usage: curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/uninstall.sh | bash

JVM_DIR="${JVM_DIR:-$HOME/.jvm}"

echo "⚠️  This will uninstall jvm from your system."

# Helper to read prompt interactively even when piped via curl | bash
prompt_user() {
    local prompt_msg="$1"
    local result=""
    if [ -t 0 ]; then
        read -r -p "$prompt_msg" result
    elif [ -e /dev/tty ]; then
        read -r -p "$prompt_msg" result < /dev/tty
    else
        result="n"
    fi
    echo "$result"
}

# 1. Main confirmation prompt
CONFIRM=$(prompt_user "Are you sure you want to uninstall jvm? [y/N] ")
case "$CONFIRM" in
    [yY][eE][sS]|[yY])
        ;;
    *)
        echo "Uninstallation cancelled."
        exit 0
        ;;
esac

# 2. Nested prompt: ask whether to keep or remove downloaded JDKs
KEEP_JDKS="y"
if [ -d "$JVM_DIR/versions" ] && [ "$(ls -A "$JVM_DIR/versions" 2>/dev/null)" ]; then
    echo ""
    echo "Found installed JDK versions in $JVM_DIR/versions."
    JDK_CHOICE=$(prompt_user "Do you want to keep your installed JDKs? [Y/n] ")
    case "$JDK_CHOICE" in
        [nN][oO]|[nN])
            KEEP_JDKS="n"
            ;;
        *)
            KEEP_JDKS="y"
            ;;
    esac
fi

echo ""
echo "Uninstalling jvm..."

# 3. Clean up shell profiles
clean_profile() {
    local profile="$1"
    if [ -f "$profile" ] && grep -q "jvm" "$profile"; then
        echo "Removing jvm configuration from $profile..."
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' '/# >>> jvm initialization >>>/,/# <<< jvm initialization <<</d' "$profile" 2>/dev/null || true
            sed -i '' '/jvm\.sh/d' "$profile" 2>/dev/null || true
        else
            sed -i '/# >>> jvm initialization >>>/,/# <<< jvm initialization <<</d' "$profile" 2>/dev/null || true
            sed -i '/jvm\.sh/d' "$profile" 2>/dev/null || true
        fi
        echo "✓ Cleaned $profile"
    fi
}

clean_profile "$HOME/.zshrc"
clean_profile "$HOME/.bashrc"
clean_profile "$HOME/.bash_profile"
clean_profile "$HOME/.config/fish/config.fish"

# 4. Remove jvm files according to choice
if [ -d "$JVM_DIR" ]; then
    if [ "$KEEP_JDKS" = "y" ]; then
        echo "Keeping installed JDKs in $JVM_DIR/versions."
        rm -rf "$JVM_DIR/bin" "$JVM_DIR/downloads" "$JVM_DIR/current" "$JVM_DIR/jvm.sh" "$JVM_DIR/jvm.fish" "$JVM_DIR/config.json"
        echo "✓ Removed jvm binary, scripts, and runtime files"
    else
        echo "Removing all jvm files and JDK versions from $JVM_DIR..."
        rm -rf "$JVM_DIR"
        echo "✓ Removed $JVM_DIR"
    fi
fi

echo ""
echo "🎉 jvm was successfully uninstalled!"
echo "Please restart your terminal session to complete."
echo ""
