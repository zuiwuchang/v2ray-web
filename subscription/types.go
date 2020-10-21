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
	Network           string             `json:"network,omitempty"`
	Security          string             `json:"security,omitempty"`
	WebsocketSettings *WebsocketSettings `json:"wsSettings,omitempty"`
	HTTPSettings      *HTTPSettings      `json:"httpSettings,omitempty"`
	TLSSettings       *TLSSettings       `json:"tlsSettings,omitempty"`
}

// WebsocketSettings .
type WebsocketSettings struct {
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// HTTPSettings .
type HTTPSettings struct {
	Path string   `json:"path,omitempty"`
	Host []string `json:"host,omitempty"`
}

// TLSSettings .
type TLSSettings struct {
	AllowInsecure bool `json:"allowInsecure,omitempty"`
}
