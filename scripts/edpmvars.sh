#!/bin/bash

usage () {
    echo "Usage: $0 [-h] [-v] [-d] command <args>"
    echo "  -h  Display this help message"
    echo "  -v  Enable verbose mode"
    echo "  -d  Enable debug mode"
    echo "  extract edpm_repo_dir                   -- Extract edpm variables from edpm repository"
    echo "  match edpm_var_file config_download_dir -- Find edpm variables defined in edpm_var_file"
    echo "                                             in the config-download directory"
}

extract_edpm_vars() {
    local edpm_dir=$1
    
    TEMPFILE=$(mktemp --suffix .yaml)
    
    while IFS= read -r -d '' file
    do
        (( count++ ))
        echo "Searching $file"
        yq '.argument_specs.main.options[]|key' "$file" | sed 's/edpm_//g'>> "$TEMPFILE"
    done <   <(find "$edpm_dir" -name 'argument_specs.yml' -print0)
    echo "Searched $count files, output in $TEMPFILE"
}

match_edpm_vars() {
    local edpm_var_file=$1
    local config_dir=$2
    
    TEMPFILE="/tmp/varmap.yaml"
    
    # Read $edpm_var_file line by line into an array
    mapfile -t lines < "$edpm_var_file"
    
    mapped=0
    totalVars=${#lines[@]}
    echo "#" > "$TEMPFILE"
    # Iterate over the array
    for line in "${lines[@]}"; do
        echo -n "Mapping $line: "
        mapped_line="${line/edpm_/tripleo_}"
        echo -n "$line: " >> "$TEMPFILE"
        
        #find ../overcloud-deploy -name '*.yaml' -exec sed -n -E 's/^[[:space:]]*(tripleo_.*):.*/\1/p' {} \;
        readarray -t files < <(find "$config_dir" -name '*.yaml' -print0 | xargs -0 sed -n -E 's/^[[:space:]]*('"$mapped_line"'):.*/\1/p')
        
        if [ ${#files[@]} -gt 0 ]; then
            echo "$mapped_line" >> "$TEMPFILE"
            echo "mapped"
            mapped=$((mapped+1))
        else
            echo " \"\"" >> "$TEMPFILE"
            echo "none"
        fi
        
    done
    echo "Mapped $mapped/$totalVars variables..."
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

if [ "$#" -lt 2 ]; then
    usage
    exit 0
fi

COMMAND=$1
shift

case "$COMMAND" in
    # Parse options to the install sub command
    extract)
        if [ "$#" -lt 1 ]; then
            usage
            exit 1
        fi
        edpm_repo=$(realpath "$1")
        extract_edpm_vars "$edpm_repo" "$2"
    ;;
    match)
        if [ "$#" -lt 2 ]; then
            usage
            exit 1
        fi
        match_edpm_vars "$1" "$2"
    ;;
    *)
        echo "Unknown command: $COMMAND"
        usage
        exit 1
    ;;
esac
