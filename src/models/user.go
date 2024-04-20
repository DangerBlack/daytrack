package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	Salt      string `json:"salt"`
}
