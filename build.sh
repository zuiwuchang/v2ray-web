#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/script/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "build script"
    echo
    echo "Usage:"
    echo "  $0 [flags]"
    echo "  $0 [command]"
    echo
    echo "Available Commands:"
    echo "  help              help for $0"
    echo "  clear             clear output"
    echo "  go                go build helper"
    echo "  docker            go docker helper"
    if [[ $View == 1 ]];then
        echo "  view              view build helper"
    fi
    echo "  pack              pack release"
    echo "  run               run project"
    echo
    echo "Flags:"
    echo "  -e, --export        test export"
    echo "  -i, --import        test import"
    echo "  -h, --help          help for $0"
}

function doexport
{
    cd "$BashDir/bin"
    if [[ ! -d temp ]];then
        mkdir temp
    fi
    echo 'v2ray-web export \'
    echo '    -s temp/settings.json \'
    echo '    --strategy temp/strategy.json\' 
    echo '    -v temp/v2ray.txt \'
    echo '    --subscription temp/subscription.json \'
    echo '    --iptables temp/iptables.json \'
    echo '    --iptables-view temp/iptables-view.txt \'
    echo '    --iptables-clear temp/iptables-clear.txt \'
    echo '    --iptables-init temp/iptables-init.txt \'
    echo '    -u temp/user.json \'
    echo '    -p temp/proxy.json'
    ./v2ray-web export \
        -s temp/settings.json \
        --strategy temp/strategy.json\
        -v temp/v2ray.txt \
        --subscription temp/subscription.json \
        --iptables temp/iptables.json \
        --iptables-view temp/iptables-view.txt \
        --iptables-clear temp/iptables-clear.txt \
        --iptables-init temp/iptables-init.txt \
        -u temp/user.json \
        -p temp/proxy.json
}
function doimport
{
    cd "$BashDir/bin"
    echo 'v2ray-web import \'
    echo '    -s temp/settings.json \'
    echo '    --strategy temp/strategy.json\' 
    echo '    -v temp/v2ray.txt \'
    echo '    --subscription temp/subscription.json \'
    echo '    --iptables temp/iptables.json \'
    echo '    --iptables-view temp/iptables-view.txt \'
    echo '    --iptables-clear temp/iptables-clear.txt \'
    echo '    --iptables-init temp/iptables-init.txt \'
    echo '    -u temp/user.json \'
    echo '    -p temp/proxy.json'
    ./v2ray-web import \
        -s temp/settings.json \
        --strategy temp/strategy.json\
        -v temp/v2ray.txt \
        --subscription temp/subscription.json \
        --iptables temp/iptables.json \
        --iptables-view temp/iptables-view.txt \
        --iptables-clear temp/iptables-clear.txt \
        --iptables-init temp/iptables-init.txt \
        -u temp/user.json \
        -p temp/proxy.json
}
case "$1" in
    help|-h|--help)
        help
    ;;
    -e|--export)
        
        doexport
        exit 0
    ;;
    -i|--import)
        
        doimport
        exit 0
    ;;
    clear)
        shift
        export Command="$0 clear"
        "$BashDir/script/clear.sh" "$@"
    ;;
    pack)
        shift
        export Command="$0 pack"
        "$BashDir/script/pack.sh" "$@"
    ;;
    go)
        shift
        export Command="$0 go"
        "$BashDir/script/go.sh" "$@"
    ;;
    docker)
        shift
        export Command="$0 docker"
        "$BashDir/script/docker.sh" "$@"
    ;;
    view)
        if [[ $View == 1 ]];then
            shift
            export Command="$0 view"
            "$BashDir/script/view.sh" "$@"
        else
            echo Error: unknown command "$1" for "$0"
            echo "Run '$0 --help' for usage."
            exit 1
        fi
    ;;
    run)
        shift
        export Command="$0 run"
        "$BashDir/script/run.sh" "$@"
    ;;
    *)
        if [[ "$1" == "" ]];then
            help
        elif [[ "$1" == -* ]];then
            echo Error: unknown flag "$1" for "$0"
            echo "Run '$0 --help' for usage."
        else
            echo Error: unknown command "$1" for "$0"
            echo "Run '$0 --help' for usage."
        fi        
        exit 1
    ;;
esac