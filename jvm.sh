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
    # Check if a .java-version or .jvmrc exists anywhere in this directory or its parents
    local check_dir="$PWD"
    local found=""
    while [ "$check_dir" != "/" ] && [ "$check_dir" != "." ]; do
        if [ -f "$check_dir/.java-version" ] || [ -f "$check_dir/.jvmrc" ]; then
            found="1"
            break
        fi
        check_dir="$(dirname "$check_dir")"
    done

    if [ -n "$found" ]; then
        eval "$(command jvm env --shell posix)" 2>/dev/null
    fi
}

if [ -n "$ZSH_VERSION" ]; then
    autoload -U add-zsh-hook 2>/dev/null && add-zsh-hook chpwd _jvm_auto_switch 2>/dev/null
elif [ -n "$BASH_VERSION" ]; then
    # Hook for Bash PROMPT_COMMAND
    if [[ ";$PROMPT_COMMAND;" != *";_jvm_auto_switch;"* ]]; then
        PROMPT_COMMAND="_jvm_auto_switch${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
    fi
fi

# Initialize active or default java on shell startup
eval "$(command jvm env --shell posix)" 2>/dev/null
