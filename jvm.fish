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

command jvm env --shell fish | source 2>/dev/null
