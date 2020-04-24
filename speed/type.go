package speed

import (
	"gitlab.com/king011/v2ray-web/db/data"
)

// StatusRunning 正在運行
const StatusRunning = 1

// StatusError 錯誤
const StatusError = 2

// StatusOk 完成
const StatusOk = 3

// Element .
type Element struct {
	ID       uint64
	Outbound data.Outbound
}

// Result 響應數據
type Result struct {
	Status   int    `json:"status,omitempty"`
	ID       uint64 `json:"id,omitempty"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration,omitempty"`
}
