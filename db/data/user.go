package data

// UserBucket .
const UserBucket = "user"

// User 用戶
type User struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:"yaml,omitempty"`
	Password string `json:"-" xml:"-" yaml:"-"`
}
