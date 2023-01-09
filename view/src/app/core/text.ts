export const ContextText = `{
    BasePath string      // 程式運行路徑
    AddIP string         // 服務器 ip
    Outbound {
      Name string        // 方案名稱
      Add string         // 服務器地址
      Port string        // 服務器端口
      Host string        // 服務器主機名稱
      TLS string         // tls 方案
      Net string         // 傳輸層協議
      Path string        // 請求路由
      UserID string      // 用戶 id 或 密碼
      AlterID string     // 用戶標識
      Security string    // 數據加密方案
      Level string       // 用戶等級
      Protocol string    // 代理協議
      Flow string        // xtls 流控
    }
  }`
