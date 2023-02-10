#!/bin/bash
set -e
if [[ "$@" == "default-command" ]];then
    # import init settings
    if [[ ! -f  /data/v2ray-web.db ]]; then
        if [[ "$V2RAY_SETTINGS" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db -s "$V2RAY_SETTINGS"
        fi
        if [[ "$V2RAY_STRATEGY" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --strategy "$V2RAY_STRATEGY"
        fi
        if [[ "$V2RAY_V2RAY" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db -v "$V2RAY_V2RAY"
        fi
        if [[ "$V2RAY_SUBSCRIPTION" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --subscription "$V2RAY_SUBSCRIPTION"
        fi
        if [[ "$V2RAY_IPTABLES" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --iptables "$V2RAY_IPTABLES"
        fi
        if [[ "$V2RAY_IPTABLES_VIEW" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --iptables-view "$V2RAY_IPTABLES_VIEW"
        fi
        if [[ "$V2RAY_IPTABLES_CLEAR" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --iptables-clear "$V2RAY_IPTABLES_CLEAR"
        fi
        if [[ "$V2RAY_IPTABLES_INIT" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db --iptables-init "$V2RAY_IPTABLES_INIT"
        fi
        if [[ "$V2RAY_USER" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db -u "$V2RAY_USER"
        fi
        if [[ "$V2RAY_PROXY" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db -p "$V2RAY_PROXY"
        fi
        if [[ "$V2RAY_LAST" != "" ]]; then
            /opt/v2ray-web/v2ray-web import --db /data/v2ray-web.db -l "$V2RAY_LAST"
        fi
    fi
    if [[ "$V2RAY_ADDR" == "" ]];then
        exec /opt/v2ray-web/v2ray-web web --no-upgrade
    else
        exec /opt/v2ray-web/v2ray-web web --no-upgrade -a "$V2RAY_ADDR"
    fi
else
    exec "$@"
fi
