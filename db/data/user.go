package data

// UserBucket .
const UserBucket = "user"

// User 用戶
type User struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	Password string `json:"-" xml:"-" yaml:"-"`
}

type UserRaw struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	Password string `json:"password,omitempty" xml:"password,omitempty" yaml:"password,omitempty"`
}
