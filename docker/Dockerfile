FROM alpine:3.17.1

RUN apk add --no-cache \
    iptables iproute2 bash

# 拷貝 v2ray-web
Copy cnf /cnf

RUN set -eux; \
    mkdir /opt/v2ray-web; \
    tar -zxvf /cnf/bin.tar.gz -C /opt/v2ray-web; \
    cp /cnf/v2ray-web.jsonnet /opt/v2ray-web/; \
    mv /cnf/proxychains.conf /etc/proxychains.conf; \
    cp /cnf/v2ray-web /opt/v2ray-web/v2ray-web; \
    rm /cnf -rf; 

Copy docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
VOLUME /data
EXPOSE 1080
EXPOSE 8118
EXPOSE 80
CMD ["default-command"]
