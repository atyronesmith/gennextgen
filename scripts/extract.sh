#!/bin/bash
# ssh -L 8080:192.168.122.120:80 nextgen

INPUT_STACK=.

extract_network_info() {
    find "$INPUT_STACK"
}

dummy() {
    find . -type f -exec grep -Il . {} + -print0 | xargs sed -n -r 's/.*\b([a-z][0-9a-z_]+).*/\1/p' | sort | uniq
    oslo-config-generator --format json --config-file config-generator/compute.conf | jq
}

VERBOSE="false"
export VERBOSE

while getopts ":hvo:di:" opt; do
    case ${opt} in
    o)
        out_dir=$OPTARG
        ;;
    i)
        INPUT_STACK=$OPTARG
        ;;
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

if [ "$#" -gt 0 ]; then
    COMMAND=$1
    shift
else
    usage
fi

case "$COMMAND" in
# Parse options to the install sub command
deploy)
    if [ "$#" -lt 1 ]; then
        usage
    fi
    deploy "$@"
    ;;
destroy)
    if [ "$#" -lt 1 ]; then
        usage
    fi
    destroy "$@"
    ;;
create-deploy)
    create_deploy
    ;;
prep-osp)
    prepare_openstack
    ;;
test-sriov)
    test_sriov
    ;;
patch-ocp)
    patch_ocp
    ;;
prep-ocp)
    prepare_for_ocp_worker
    ;;
csr)
    sign_csr
    ;;
clean)
    rm -rf "$BUILD_DIR"
    ;;
label-nodes)
    label_nodes feature.node.kubernetes.io/network-sriov.capable="true"
    ;;
deploy-operator)
    if [ "$#" -lt 1 ]; then
        usage
    fi
    deploy_operator "$1"
    ;;
deploy-policy)
    deploy_policy
    ;;
pull-secret)
    if [ "$#" -lt 1 ]; then
        usage
    fi
    create_pull_secret "$1"
    ;;
install-manifests)
    install_manifests
    ;;
ingress-fip)
    create_ingress_fip
    ;;
*)
    echo "Unknown command: $COMMAND"
    usage "$out_dir"
    ;;
esac
