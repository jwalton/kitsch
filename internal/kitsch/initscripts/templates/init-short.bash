__main() {
    local major="${BASH_VERSINFO[0]}"
    local minor="${BASH_VERSINFO[1]}"

    if ((major > 4)) || { ((major == 4)) && ((minor >= 1)); }; then
        source <({{ .kitschCommand }} init {{with .configFile}}--config {{.}} {{end}}--print-full-init bash)
    else
        source /dev/stdin <<<"$({{ .kitschCommand }} init {{with .configFile}}--config {{.}} {{end}}--print-full-init bash)"
    fi
}
__main
unset -f __main