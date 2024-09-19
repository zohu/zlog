#!/bin/bash

args=()
while [ $# -gt 0 ]; do
    case "$1" in
        --help)
            echo "Usage: $0 [options]"
            exit 0
            ;;
        --*)
            log_error_with_exit "unknown option: $1"
            ;;
        *)
            args+=("$1")
            ;;
    esac
    shift $(( $# > 0 ? 1 : 0 ))
done

log_info() {
    printf "[INFO] %s\n" "$*"
}
log_prompt() {
    printf "\033[0;34m%s\033[0m\n" "$*"
}
log_error() {
    printf "\033[0;31m[ERROR] %s\033[0m\n" "$*"
}
log_error_exit() {
    log_error "$*"
    exit 1
}

cmd_exists() {
    command -v "$*" > /dev/null 2>&1
}
cmd() {
    log_prompt "+ $*"
    /bin/bash -c "$*"
}
install() {
    if [ -z "$1" ]; then
        log_error_exit "URL is required, build.sh go-install <URL> [options...]"
    fi
    name=${1##*/}
    name=${name%%@*}
    if ! cmd_exists "$name"; then
        cmd go install "$1"
    fi
}
protoc(){
    cmd "find . -type f -name '*.pb.go' -exec rm -rf {} \;"
    cmd "find . -type f -name '*.proto' -exec protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. {} \;"
    cmd "find . -type f -name '*.pb.go' -exec protoc-go-inject-tag -input={} \;"
    cmd "go mod tidy"
}

main(){
    if [ ${#args[*]} -lt 1 ]; then
        log_error_exit "Usage: $0 <action> [args...] [options...]"
    fi

    case "${args[0]}" in
        install)
            install "${args[*]:1}"
            ;;
        protoc)
            protoc
            ;;
        *)
            log_error_exit "unknown action: ${args[0]}"
            ;;
    esac
}
main