# v2ray-web

v2ray-web 是 v2ray-core 的一個 web接口的 ui程序 爲桌面系統提供了一個 簡單且一致的 操作方案

# 安裝與運行

[下載](https://gitlab.com/king011/v2ray-web/-/releases) 各平臺對應的 壓縮包 並解壓

執行 v2ray-web web 指令 運行程式

打開瀏覽器 訪問 http://127.0.0.1:1989 控制v2ray-core


# 預覽
![](document/view.png)
![](document/about.png)
![](document/template.png)

# docker

docker 容器打包了 v2ray-web 以及一些網路測試與代理相關的工具，你可以使用它來運行一個 xray 的代理程序。


```
docker run \
    --name v2ray-web \
    -p 8080:80/tcp \
    -p 1080:1080/tcp \
    -p 28118:8118/tcp \
    -v /data:YourConfigDir  \
    -d king011/v2ray-web:v1.6.0
```

* 8080 端口用於訪問 v2ray-web 提供的網頁 ui
* 1080 是 socks5 代理
* 8118 是 http 代理

容器裏面也已經打包了 iptables，使用它可以在容器裏面配置透明代理，但需要指定 `--cap-add=NET_ADMIN --cap-add=NET_RAW`

```
docker run \
    --name v2ray-web \
    --cap-add=NET_ADMIN --cap-add=NET_RAW \
    -p 8080:80/tcp \
    -p 1080:1080/tcp \
    -p 28118:8118/tcp \
    -v /data:YourConfigDir  \
    -d king011/v2ray-web:v1.6.0
```

通常不建議使用這個，除非你剛好在朝鮮之類的地方，那麼你可以使用這個鏡像作爲基礎鏡像，然後構建一個支持透明代理的 dev 容器，這樣你可以在 dev 容器裏面使用透明代理進行開發工作，因爲容器與宿主機有各自獨立的網卡所以 dev 容器的透明代理不會影響宿主機的網路