export const ContextText = `{
    BasePath string         // 程式運行路徑
    AddIP string            // 服務器 ip
    Outbound {
      Name string           // 方案名稱
      Add string            // 服務器地址
      Port string           // 服務器端口
      Host string           // 服務器主機名稱
      TLS string            // tls 方案
      Net string            // 傳輸層協議
      Path string           // 請求路由
      UserID string         // 用戶 id 或 密碼
      AlterID string        // 用戶標識
      Security string       // 數據加密方案
      Level string          // 用戶等級
      Protocol string       // 代理協議
      Flow string           // xtls 流控
    }
    Strategy {
      Name string           // 策略名稱
      Value int             // 策略值
      Host [][]string       // 靜態 ip 列表 [['baidu.com', '127.0.0.1'], ['dns.google', '8.8.8.8', '8.8.4.4']]
      ProxyIP []string      // 這些 ip 使用代理
      ProxyDomain []string  // 這些 域名 使用代理
      DirectIP []string     // 這些 ip 直連
      DirectDomain []string // 這些 域名 直連
      BlockIP []string      // 這些 ip 阻止訪問
      BlockDomain []string  // 這些 域名 阻止訪問
    }
  }`
