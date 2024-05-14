#!/bin/bash

usage () {
    echo "Usage: $0 [-h] [-v] [-d] va_dir"
    echo "  -h  Display this help message"
    echo "  -v  Enable verbose mode"
    echo "  -d  Enable debug mode"
    echo "  va_dir path to VA repository"
}

extract_arch_vars() {
    local arch_dir=$1
    
    TEMPFILE=$(mktemp --suffix .yaml)
    
    while IFS= read -r -d '' file
    do
        (( count++ ))
        echo "Searching $file"
        yq '.replaces' "$file"
    done <   <(find "$arch_dir" -name 'kustomization.yaml' -print0)
    echo "Searched $count files, output in $TEMPFILE"
}

VERBOSE="false"
export VERBOSE

while getopts "hvd" opt; do
    case ${opt} in
        v)
            VERBOSE="true"
        ;;
        d)
            set -x
        ;;
        h)
            usage
            exit 0
        ;;
        \?)
            echo "Invalid Option: -$OPTARG" 1>&2
            exit 1
        ;;
    esac
done

shift $((OPTIND - 1))

if [ "$#" -ne 1 ]; then
    usage
fi
arch_dir=$(realpath "$1")

extract_arch_vars "$arch_dir"