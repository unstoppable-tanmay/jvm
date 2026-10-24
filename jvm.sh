# jvm shell integration (sourced in ~/.zshrc, ~/.bashrc, etc.)
# Similar to nvm.sh

export JVM_DIR="${JVM_DIR:-$HOME/.jvm}"

# Ensure jvm binary is in PATH
if [ -d "$JVM_DIR/bin" ] && [[ ":$PATH:" != *":$JVM_DIR/bin:"* ]]; then
    export PATH="$JVM_DIR/bin:$PATH"
fi

# Thin shell wrapper function
jvm() {
    if [ "$1" = "use" ]; then
        command jvm "$@" || return $?
        eval "$(command jvm env --shell posix)"
    else
        command jvm "$@"
    fi
}

# Auto-switch Java version when cd-ing into a directory with .java-version or .jvmrc
_jvm_auto_switch() {
    if [ -f .java-version ] || [ -f .jvmrc ]; then
        eval "$(command jvm env --shell posix)" 2>/dev/null
    fi
}

if [ -n "$ZSH_VERSION" ]; then
    autoload -U add-zsh-hook 2>/dev/null && add-zsh-hook chpwd _jvm_auto_switch 2>/dev/null
fi

# Initialize default java on shell startup
eval "$(command jvm env --shell posix)" 2>/dev/null
