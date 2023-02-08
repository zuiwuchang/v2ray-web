Target="v2ray-web"
Docker="king011/v2ray-web"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v1.7.2"
View=1
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)