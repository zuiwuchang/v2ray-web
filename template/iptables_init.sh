#!/bin/bash
set -e

PROXY_PORT=12345

Whitelist=(
    # 私有地址
	0.0.0.0/8
	127.0.0.0/8
	10.0.0.0/8
    172.0.0.0/8
	169.0.0.0/8
	192.0.0.0/8
    # 廣播地址
	224.0.0.0/4
	240.0.0.0/4
)

# 添加路由表 100
if [[ `ip rule list  | egrep '0x1 lookup 100'` == "" ]];then
    ip rule add fwmark 1 table 100
fi
# 爲路由表 100 設定規則
if [[ `ip route list table 100 | egrep 'dev lo'` == "" ]];then
    ip route add local 0.0.0.0/0 dev lo table 100
fi

# 已經設置過直接返回
if [[ `iptables-save |egrep '\-A OUTPUT \-p udp \-j XRAY_SELF'` != "" ]];then
    exit 0
fi

# 創建鏈
iptables -t mangle -N XRAY
iptables -t mangle -N XRAY_SELF
iptables -t mangle -N XRAY_DIVERT

# 放行私有地址與廣播地址
for whitelist in "${Whitelist[@]}"
do
        iptables -t mangle -A XRAY -d "$whitelist" -j RETURN
        iptables -t mangle -A XRAY_SELF -d "$whitelist" -j RETURN
done

# 可選的配置避免已有連接的包二次通過 tproxy 從而提升一些性能
iptables -t mangle -A XRAY_DIVERT -j MARK --set-mark 1
iptables -t mangle -A XRAY_DIVERT -j ACCEPT
iptables -t mangle -I PREROUTING -p tcp -m socket -j XRAY_DIVERT

# 代理局域網設備 
iptables -t mangle -A XRAY -p tcp -j TPROXY --on-port "$PROXY_PORT" --tproxy-mark 1 # tcp 到 tproxy 代理端口
iptables -t mangle -A XRAY -p udp -j TPROXY --on-port "$PROXY_PORT" --tproxy-mark 1 # udp 到 tproxy 代理端口
iptables -t mangle -A PREROUTING -j XRAY # 流量都重定向到 XRAY 鏈

# 代理本機
iptables -t mangle -A XRAY_SELF -m mark --mark 2 -j RETURN # 放行所有 mark 2 的流量
iptables -t mangle -A XRAY_SELF -j MARK --set-mark 1 # 爲流量設置 mark 1
iptables -t mangle -A OUTPUT -p tcp -j XRAY_SELF # tcp 到 XRAY_SELF 鏈
iptables -t mangle -A OUTPUT -p udp -j XRAY_SELF # udp 到 XRAY_SELF 鏈
