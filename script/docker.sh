#!/usr/bin/env bash
set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "docker build helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -r, --run           run docker"
    echo "  -p, --push           push to hub"
    echo "  -h, --help          help for $Command"
}


ARGS=`getopt -o hrp --long help,run,push -n "$Command" -- "$@"`
eval set -- "${ARGS}"
go=0
push=0
run=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -r|--run)
            run=1
            shift
        ;;
        -g|--go)
            go=1
            shift
        ;;
        -p|--push)
            push=1
            shift
        ;;
        --)
            shift
            break
        ;;
        *)
            echo Error: unknown flag "$1" for "$Command"
            echo "Run '$Command --help' for usage."
            exit 1
        ;;
    esac
done

if [[ "$run" == 1 ]];then
    args=(
        sudo docker run --rm --cap-add=NET_ADMIN --cap-add=NET_RAW -p 20080:80/tcp -p 21080:1080/tcp -p 28118:8118/tcp --name v2ray-web -d "\"$Docker:$Version\""
    )
    exec="${args[@]}"
    echo $exec
    eval "$exec"
    exit $?
fi


cd "$Dir/docker"
cp ../bin/linux.amd64.tar.gz ./cnf/bin.tar.gz
args=(
    sudo docker build --network host -t "\"$Docker:$Version\"" .
)
exec="${args[@]}"
echo $exec
eval "$exec"

if [[ "$push" == 1 ]];then
    args=(
        sudo docker push "\"$Docker:$Version\""
    )
    exec="${args[@]}"
    echo $exec
    eval "$exec"
fi
