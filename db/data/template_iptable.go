package data

import "gitlab.com/king011/v2ray-web/version"

const iptablesInit = `# tag =` + version.Tag + `
# date ` + version.Date + `
# commit ` + version.Commit + `

# 本地 dns 端口
DNS_Port=10054

# Redir 程式 端口
Redir_Port=10090

# 要放行的 ip 數組 
# 通常是 服務器地址 和 不需要代理的地址
IP_Servers=(
    {{.AddIP}}
    114.114.114.114
)

# 定義 內網 地址
# 一般不用修改
IP_Private=(
	0/8
	127/8
	10/8
	169.254/16
	172.16/12
	192.168/16
	224/4
	240/4
)

# 創建 nat/tcp 轉發鏈 用於 轉發 tcp流
iptables-save | egrep "^\:NAT_TCP" >> /dev/null
if [[ $? != 0 ]];then
    iptables -t nat -N NAT_TCP
fi

# 放行所有 內網地址
iptables -t nat -A NAT_TCP -p tcp -j RETURN -m mark --mark 0xff
for i in ${!IP_Private[@]}
do
    iptables -t nat -A NAT_TCP -d ${IP_Private[i]} -j RETURN
done

# 放行 發往 服務器的 數據
for i in ${!IP_Servers[@]}
do
    iptables -t nat -A NAT_TCP -d ${IP_Servers[i]} -j RETURN
done

# 重定向 tcp 數據包到 redir 監聽端口
iptables -t nat -A NAT_TCP -p tcp -j REDIRECT --to-ports $Redir_Port


# 重定向 向網關發送的 dns 查詢
for i in ${!IP_Private[@]}
do
	iptables -t nat -A OUTPUT -d ${IP_Private[i]} -p udp -m udp --dport 53 -j DNAT --to-destination 127.0.0.1:$DNS_Port
	iptables -t nat -A OUTPUT -d ${IP_Private[i]} -p tcp -m tcp --dport 53 -j DNAT --to-destination 127.0.0.1:$DNS_Port
done

# 重定向 數據流向 NAT_TCP
iptables -t nat -A OUTPUT -p tcp -j NAT_TCP
iptables -t nat -A PREROUTING -p tcp -s 192.168/16 -j NAT_TCP
iptables -t nat -A POSTROUTING -s 192.168/16 -j MASQUERADE
`
