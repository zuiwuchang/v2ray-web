#!/bin/bash

if [[ "$@" == "default-command" ]];then
    exec /opt/v2ray-web/v2ray-web web --no-upgrade
else
    exec "$@"
fi
