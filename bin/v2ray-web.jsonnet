local Millisecond = 1;
local Second = 1000 * Millisecond;
local Minute = 60 * Second;
local Hour = 60 * Minute;
local Day = 24 * Hour;
{
    HTTP: {
        Addr: "localhost:8080",
        // x509 if empty use http else use https
		// CertFile: "test.pem",
		// KeyFile: "test.key",
        // cookie 過期時間
        Maxage: Day * 30,
        // cookie 密鑰
        Secret: "cerberus is an idae",
		// ui界面目錄
		View:"view",
    },
	Database:{
		// 數據源 位置
		Source: "v2ray-web.db",
	},
	Logger: {
		// zap http
		//HTTP:"localhost:20000",
		// log name
		//Filename:"logs/v2ray-web.log",
		// MB
		MaxSize:    100, 
		// number of files
		MaxBackups: 3,
		// day
		MaxAge:     28,
		// level : debug info warn error dpanic panic fatal
		Level :"debug",
		// 是否要 輸出 代碼位置
        Caller:true,
	},
}