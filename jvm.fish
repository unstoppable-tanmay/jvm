# jvm fish shell integration
set -q JVM_DIR; or set -gx JVM_DIR "$HOME/.jvm"

if test -d "$JVM_DIR/bin"; and not contains "$JVM_DIR/bin" $PATH
    set -gx PATH "$JVM_DIR/bin" $PATH
end

function jvm
    if test (count $argv) -gt 0; and test $argv[1] = "use"
        command jvm $argv
        and eval (command jvm env --shell fish)
    else
        command jvm $argv
    end
end

# Auto-switch hook on directory change in fish
function _jvm_auto_switch --on-variable PWD
    set -l check_dir $PWD
    set -l found ""
    while test "$check_dir" != "/" -a "$check_dir" != "."
        if test -f "$check_dir/.java-version" -o -f "$check_dir/.jvmrc"
            set found "1"
            break
        end
        set check_dir (dirname "$check_dir")
    end

    if test -n "$found"
        command jvm env --shell fish | source 2>/dev/null
    end
end

command jvm env --shell fish | source 2>/dev/null
