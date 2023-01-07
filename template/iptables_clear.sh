#!/bin/bash
set -e

# 清空 XRAY_DIVERT
if [[ `iptables-save | egrep '\-A PREROUTING \-p tcp \-m socket \-j XRAY_DIVERT'` != "" ]];then
    iptables -t mangle -D PREROUTING -p tcp -m socket -j XRAY_DIVERT
fi
if [[ `iptables-save | egrep XRAY_DIVERT` != "" ]];then
    iptables -t mangle -F XRAY_DIVERT
    iptables -t mangle -X XRAY_DIVERT
fi

# 清空 XRAY_SELF 
if [[ `iptables-save | egrep '\-A OUTPUT \-p tcp \-j XRAY_SELF'` != "" ]];then
    iptables -t mangle -D OUTPUT -p tcp -j XRAY_SELF
fi
if [[ `iptables-save | egrep '\-A OUTPUT \-p udp \-j XRAY_SELF'` != "" ]];then
    iptables -t mangle -D OUTPUT -p udp -j XRAY_SELF
fi
if [[ `iptables-save | egrep XRAY_SELF` != "" ]];then
    iptables -t mangle -F XRAY_SELF
    iptables -t mangle -X XRAY_SELF
fi

# 清空 XRAY 
if [[ `iptables-save | egrep '\-A PREROUTING \-j XRAY'` != "" ]];then
    iptables -t mangle -D PREROUTING -j XRAY
fi
if [[ `iptables-save | egrep XRAY` != "" ]];then
    iptables -t mangle -F XRAY
    iptables -t mangle -X XRAY
fi

# 刪除 路由規則
if [[ `ip route list table 100 | egrep 'dev lo'` != "" ]];then
    ip route del local 0.0.0.0/0 dev lo table 100
fi
# 刪除 路由表
if [[ `ip rule list  | egrep '0x1 lookup 100'` != "" ]];then
    ip rule del fwmark 1 table 100
fi
