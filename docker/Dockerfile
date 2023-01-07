FROM ubuntu:22.04

# 安裝一些額外工具
ENV DEBIAN_FRONTEND=noninteractive
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends iptables iproute2 net-tools iputils-ping dnsutils ca-certificates; \
    apt-get install -y --no-install-recommends curl wget vim proxychains; \
    rm -rf /var/lib/apt/lists/*;

# 拷貝 v2ray-web
Copy cnf /cnf

RUN set -eux; \
    mkdir /opt/v2ray-web; \
    tar -zxvf /cnf/bin.tar.gz -C /opt/v2ray-web; \
    cp /cnf/v2ray-web.jsonnet /opt/v2ray-web/; \
    mv /cnf/proxychains.conf /etc/proxychains.conf; \
    rm /cnf -rf; 

Copy docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
VOLUME /data
EXPOSE 1080
EXPOSE 8118
EXPOSE 80
CMD ["default-command"]
