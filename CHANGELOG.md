# v1.6.0
* 升級 golang 到 1.19.4
* 使用 xray 1.7.1 替代 v2ray 以支持 xtls
* 支持用 es5 腳本配置 v2ray
* 升級 angular 到 14


# v1.5.0

* 升級 golang 到 1.17
* 升級 v2ray-core 到 4.44.0
* 使用 go:embed 取代 statik
* 升級 gin 到 1.7.7
* 升級 網頁前端到 angular 13
* 重寫編譯自動化腳本 build.sh

# v1.4.0
* 升級網頁前端到 angular 12
* 升級 v2ray-core 到 4.43.0
* 修復 trojan 導入遺失端口
* 修改 vless 分享格式
* 升級 gin 到 1.7.3
* 從 github 自動檢查更新

# v1.3.1
* 修復默認配置 h2 協議無法工作

# v1.3.0
* 修改默認配置檔案
* 修改默認iptables規則
* 升級 v2ray-core 到 4.36.2

# v1.0.14
* 對 1k 到 2m 的檔案 使用 gzip 壓縮 加快 網頁打開速度
* 網頁視圖部署 Service Worker
* 關於頁面 顯示 v2ray-core 版本
* 升級 v2ray-core 到 4.19.1

# v1.0.7
* 添加 說明文檔
* 修復一些bug

# v1.0.5 

使用websocket 爲網頁增加一個日誌頁面 顯示 v2ray 日誌

# v1.0.2

首次發佈 完成主要 架構和功能