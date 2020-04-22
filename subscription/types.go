package subscription

// Vnext .
type Vnext struct {
	Address string `json:"address"`
	Port    uint32 `json:"port"`
	Users   []User `json:"users"`
}

// User .
type User struct {
	ID       string `json:"id"`
	AlterID  int64  `json:"alterId"`
	Security string `json:"security"`
	Level    int64  `json:"level"`
}

// StreamSettings .
type StreamSettings struct {
	Network           string             `json:"network"`
	Security          string             `json:"security"`
	WebsocketSettings *WebsocketSettings `json:"wsSettings"`
}

// WebsocketSettings .
type WebsocketSettings struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}
